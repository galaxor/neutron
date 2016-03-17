package api

import (
	"encoding/json"

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

	sessionToken := "access_token"

	keyring := user.Addresses[0].Keys[0] // TODO: find a better way to get the keyring
	encryptedToken, err := keyring.EncryptToSelf(sessionToken)
	if err != nil {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{InternalServerError},
			Error: "invalid_key",
			ErrorDescription: err.Error(),
		})
		return
	}

	api.sessions[sessionToken] = user.ID

	ctx.JSON(200, &AuthResp{
		Resp: Resp{Ok},
		AccessToken: encryptedToken,
		ExpiresIn: 360000,
		TokenType: TokenBearer,
		Scope: "full mail payments reset keys",
		Uid: user.ID, // TODO: put something else there
		RefreshToken: "refresh_token",
		PrivateKey: user.EncPrivateKey,
		EncPrivateKey: user.EncPrivateKey,
		EventID: "event_id",
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

	userId := uid // TODO

	sessionToken := ""
	for t, id := range api.sessions {
		if id == userId {
			sessionToken = t
			break
		}
	}

	if sessionToken == "" {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{Unauthorized},
			Error: "invalid_session",
			ErrorDescription: "Not logged in",
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
		AccessToken: sessionToken,
		Uid: userId,
	})
	ctx.SetCookie("AUTH-" + sessionToken, string(authCookie), 0, "/api/", "", false, true)

	ctx.JSON(200, &AuthCookiesResp{
		Resp: Resp{Ok},
		SessionToken: sessionToken,
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
