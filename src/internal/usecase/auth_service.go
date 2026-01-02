package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pbmap_api/src/config"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/repository"
	"pbmap_api/src/pkg/auth"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	LoginWithSocial(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
}

type authService struct {
	userUsecase    UserUsecase
	tokenRepo      repository.TokenRepository
	jwtService     *auth.JWTService
	googleClientID string
	lineChannelID  string
}

func NewAuthService(userUsecase UserUsecase, tokenRepo repository.TokenRepository, jwtService *auth.JWTService, cfg *config.Config) AuthService {
	return &authService{
		userUsecase:    userUsecase,
		tokenRepo:      tokenRepo,
		jwtService:     jwtService,
		googleClientID: cfg.GoogleClientID,
		lineChannelID:  cfg.LineChannelID,
	}
}

func (s *authService) LoginWithSocial(ctx context.Context, req *dto.SocialLoginRequest) (*dto.LoginResponse, error) {
	var providerID string
	var err error

	switch req.Provider {
	case "google":
		providerID, err = s.verifyGoogleToken(req.AccessToken)
	case "line":
		providerID, err = s.verifyLineToken(req.AccessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("provider verification failed: %v", err)
	}

	user, err := s.userUsecase.SyncUserFromSocial(ctx, req.Provider, providerID)
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

	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(appTokenTTL.Seconds()),
	}, nil
}

func (s *authService) verifyGoogleToken(token string) (string, error) {
	if s.googleClientID == "" || s.googleClientID == "mock" {
		return token, nil
	}

	payload, err := idtoken.Validate(context.Background(), token, s.googleClientID)
	if err != nil {
		return "", fmt.Errorf("invalid google token: %v", err)
	}

	return payload.Subject, nil
}

func (s *authService) verifyLineToken(token string) (string, error) {
	apiURL := "https://api.line.me/oauth2/v2.1/verify"
	data := url.Values{}
	data.Set("id_token", token)
	data.Set("client_id", s.lineChannelID)

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return "", fmt.Errorf("failed to verify line id token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("line id token verification failed: status=%d, body=%s", resp.StatusCode, string(bodyBytes))
	}

	var result dto.LineIDTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode line response: %v", err)
	}

	if result.Sub == "" {
		return "", fmt.Errorf("line id token valid but sub is empty")
	}

	return result.Sub, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRepo.DeleteAppToken(ctx, userID.String())
}
