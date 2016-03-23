package api

import (
	"encoding/json"
	"strings"

	"gopkg.in/macaron.v1"
)

type GrantType string

const (
	GrantPassword GrantType = "password"
	GrantRefreshToken = "refresh_token"
)

type ResponseType string

const (
	ResponseToken ResponseType = "token"
)

type TokenType string

const (
	TokenBearer TokenType = "Bearer"
)

type AuthReq struct {
	Req
	ClientID string
	ClientSecret string
	GrantType GrantType
	Password string
	RedirectURI string
	ResponseType ResponseType
	State string
	Username string
}

type AuthResp struct {
	Resp
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

type AuthCookiesReq struct {
	Req
	ClientID string
	ResponseType ResponseType
	GrantType GrantType
	RefreshToken string
	RedirectURI string
	State string
}

type AuthCookiesResp struct {
	Resp
	SessionToken string
}

type AuthCookie struct {
	AccessToken string
	Uid string `json:"UID"`
}

func (api *Api) Auth(ctx *macaron.Context, req AuthReq) {
	if req.GrantType != GrantPassword {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{BadRequest},
			Error: "invalid_grant",
			ErrorDescription: "GrantType must be set to password",
		})
		return
	}

	user, err := api.backend.Auth(req.Username, req.Password)
	if err != nil {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{Unauthorized},
			Error: "invalid_grant",
			ErrorDescription: err.Error(),
		})
		return
	}

	session := &Session{
		ID: "session_id", // TODO: generate this
		Token: "access_token", // TODO: generate this
		UserID: user.ID,
	}

	keyring := user.GetMainAddress().Keys[0] // TODO: find a better way to get the keyring
	encryptedToken, err := keyring.EncryptToSelf(session.Token)
	if err != nil {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{InternalServerError},
			Error: "invalid_key",
			ErrorDescription: err.Error(),
		})
		return
	}

	lastEvent, err := api.backend.GetLastEvent(user.ID)
	if err != nil {
		ctx.JSON(500, newErrorResp(err))
		return
	}

	api.sessions[session.ID] = session

	ctx.JSON(200, &AuthResp{
		Resp: Resp{Ok},
		AccessToken: encryptedToken,
		ExpiresIn: 360000, // TODO: really expire
		TokenType: TokenBearer,
		Scope: "full mail payments reset keys",
		Uid: session.ID,
		RefreshToken: "refresh_token", // TODO
		PrivateKey: keyring.PrivateKey,
		EncPrivateKey: keyring.PrivateKey,
		EventID: lastEvent.ID,
	})
}

func (api *Api) AuthCookies(ctx *macaron.Context, req AuthCookiesReq) {
	uid := api.getUid(ctx)
	if uid == "" {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{BadRequest},
			Error: "invalid_grant",
			ErrorDescription: "No uid provided",
		})
		return
	}

	session, ok := api.sessions[uid]
	if !ok {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{BadRequest},
			Error: "invalid_session",
			ErrorDescription: "Invalid UID",
		})
		return
	}

	auth, ok := ctx.Req.Header["Authorization"]
	if !ok || len(auth) == 0 {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{BadRequest},
			Error: "invalid_authorization",
			ErrorDescription: "Invalid authorization header",
		})
		return
	}

	parts := strings.SplitN(auth[0], " ", 2)
	if len(parts) != 2 {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{BadRequest},
			Error: "invalid_authorization",
			ErrorDescription: "Invalid authorization header",
		})
		return
	}

	tokenType := parts[0]
	token := parts[1]

	if TokenType(tokenType) != TokenBearer || token != session.Token {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{BadRequest},
			Error: "invalid_authorization",
			ErrorDescription: "Invalid authorization header",
		})
		return
	}

	if req.GrantType != GrantRefreshToken {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{BadRequest},
			Error: "invalid_grant",
			ErrorDescription: "GrantType must be set to refresh_token",
		})
		return
	}

	authCookie, _ := json.Marshal(&AuthCookie{
		AccessToken: session.Token,
		Uid: session.ID,
	})
	ctx.SetCookie("AUTH-" + session.Token, string(authCookie), 0, "/api/", "", false, true)

	ctx.JSON(200, &AuthCookiesResp{
		Resp: Resp{Ok},
		SessionToken: session.Token,
	})
}

func (api *Api) DeleteAuth(ctx *macaron.Context) {
	sessionToken := api.getSessionToken(ctx)
	if sessionToken != "" {
		ctx.SetCookie("AUTH-" + sessionToken, "", 0, "/api/", "", false, true)

		delete(api.sessions, "sessionToken")
	}

	ctx.JSON(200, &Resp{Ok})
}
