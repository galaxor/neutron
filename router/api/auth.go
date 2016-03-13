package api

import "gopkg.in/macaron.v1"

type GrantType string

const (
	GrantPassword GrantType = "password"
)

type ResponseType string

const (
	ResponseToken ResponseType = "token"
)

type TokenType string

const (
	TokenBearer TokenType = "Bearer"
)

type AuthRequest struct {
	ClientID string
	ClientSecret string
	GrantType GrantType
	Password string
	RedirectURI string
	ResponseType ResponseType
	State string
	Username string
}

type AuthResponse struct {
	Code int
	AccessToken string
	ExpiresIn int
	TokenType TokenType
	Scope string
	Uid string
	RefreshToken string
	UserStatus int
	PrivateKey string
	EncPrivateKey string
	EventID string
}

func Auth(ctx *macaron.Context, req AuthRequest) {
	res := AuthResponse{}
	ctx.JSON(200, &res)
}
