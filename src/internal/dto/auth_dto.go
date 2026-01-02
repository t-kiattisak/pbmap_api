package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SocialLoginRequest struct {
	Provider    string `json:"provider" validate:"required,oneof=google line"` // google, line
	AccessToken string `json:"access_token" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type LineIDTokenResponse struct {
	Sub     string   `json:"sub"`
	Iss     string   `json:"iss"`
	Aud     string   `json:"aud"`
	Exp     int64    `json:"exp"`
	Iat     int64    `json:"iat"`
	Nonce   string   `json:"nonce"`
	Amr     []string `json:"amr"`
	Name    string   `json:"name"`
	Picture string   `json:"picture"`
	Email   string   `json:"email"`
}
