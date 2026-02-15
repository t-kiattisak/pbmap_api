package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"
	"pbmap_api/src/internal/dto"
	implRepositories "pbmap_api/src/internal/repositories"
	"pbmap_api/src/pkg/auth"
	"pbmap_api/src/pkg/config"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	LoginWithSocial(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userUsecase    UserUsecase
	tokenRepo      repositories.TokenRepository
	sessionRepo    repositories.SessionRepository
	tm             implRepositories.TransactionManager
	jwtService     *auth.JWTService
	googleClientID string
	lineChannelID  string
}

func NewAuthService(userUsecase UserUsecase, tokenRepo repositories.TokenRepository, sessionRepo repositories.SessionRepository, tm implRepositories.TransactionManager, jwtService *auth.JWTService, cfg *config.Config) AuthService {
	return &authService{
		userUsecase:    userUsecase,
		tokenRepo:      tokenRepo,
		sessionRepo:    sessionRepo,
		tm:             tm,
		jwtService:     jwtService,
		googleClientID: cfg.GoogleClientID,
		lineChannelID:  cfg.LineChannelID,
	}
}

func (s *authService) LoginWithSocial(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error) {
	var providerID, email, displayName string
	var err error

	switch req.Provider {
	case "google":
		providerID, email, displayName, err = s.verifyGoogleToken(req.AccessToken)
	case "line":
		providerID, email, displayName, err = s.verifyLineToken(req.AccessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("provider verification failed: %v", err)
	}

	user, err := s.userUsecase.SyncUserFromSocial(ctx, dto.CreateUserFromSocialInput{
		Provider:    req.Provider,
		ProviderID:  providerID,
		Email:       email,
		DisplayName: displayName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sync user: %v", err)
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate app token: %v", err)
	}

	appTokenTTL := 72 * time.Hour
	if err := s.tokenRepo.SetAppToken(ctx, user.ID.String(), token, appTokenTTL); err != nil {
		return nil, fmt.Errorf("failed to store app token: %v", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	err = s.tm.Do(ctx, func(ctx context.Context) error {
		deviceID := uuid.Nil
		if req.DeviceID != "" || req.DeviceType != "" || req.PushToken != "" {
			provider := "apns"
			if req.DeviceType == "android" {
				provider = "fcm"
			}

			device := &entities.UserDevice{
				UserID:     user.ID,
				Provider:   provider,
				DeviceType: req.DeviceType,
				PushToken:  req.PushToken,
				LastSeen:   time.Now(),
			}

			if err := s.userUsecase.UpsertDevice(ctx, device); err != nil {
				return fmt.Errorf("failed to register device: %v", err)
			}

			if device.ID != uuid.Nil {
				deviceID = device.ID
			}
		}

		var existingSession *entities.UserSession
		if deviceID != uuid.Nil {
			existingSession, _ = s.sessionRepo.GetSessionByDeviceID(ctx, user.ID, deviceID)
		}

		if existingSession != nil {
			existingSession.RefreshToken = refreshToken
			existingSession.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)

			if err := s.sessionRepo.UpdateSession(ctx, existingSession); err != nil {
				return fmt.Errorf("failed to update session: %v", err)
			}
		} else {
			session := &entities.UserSession{
				ID:           uuid.New(),
				UserID:       user.ID,
				RefreshToken: refreshToken,
				DeviceID:     deviceID,
				ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
			}

			if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
				return fmt.Errorf("failed to create session: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(appTokenTTL.Seconds()),
	}, nil
}

func (s *authService) verifyGoogleToken(token string) (string, string, string, error) {
	payload, err := idtoken.Validate(context.Background(), token, s.googleClientID)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid google token: %v", err)
	}

	fmt.Println(payload)

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	return payload.Subject, email, name, nil
}

func (s *authService) verifyLineToken(token string) (string, string, string, error) {
	apiURL := "https://api.line.me/oauth2/v2.1/verify"
	data := url.Values{}
	data.Set("id_token", token)
	data.Set("client_id", s.lineChannelID)

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to verify line id token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("line id token verification failed: status=%d, body=%s", resp.StatusCode, string(bodyBytes))
	}

	var result dto.LineIDTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", "", fmt.Errorf("failed to decode line response: %v", err)
	}

	if result.Sub == "" {
		return "", "", "", fmt.Errorf("line id token valid but sub is empty")
	}

	return result.Sub, result.Email, result.Name, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRepo.DeleteAppToken(ctx, userID.String())
}

func (s *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	session, err := s.sessionRepo.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessionRepo.RevokeSession(ctx, session.ID)
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.userUsecase.GetUser(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := s.sessionRepo.RevokeSession(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old session: %v", err)
	}

	newAccessToken, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	appTokenTTL := 72 * time.Hour
	if err := s.tokenRepo.SetAppToken(ctx, user.ID.String(), newAccessToken, appTokenTTL); err != nil {
		return nil, fmt.Errorf("failed to store app token: %v", err)
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	newSession := &entities.UserSession{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.sessionRepo.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create new session: %v", err)
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(appTokenTTL.Seconds()),
	}, nil
}

func (s *authService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
