package api

import (
	"bytes"
	"errors"
	"strings"
	"encoding/json"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
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

func encrypt(user *backend.User, token string) (encrypted string, err error) {
	entitiesList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(user.EncPrivateKey))
	if err != nil {
		return
	}
	if len(entitiesList) == 0 {
		err = errors.New("Key ring does not contain any key")
		return
	}

	entity := entitiesList[0]

	var tokenBuffer bytes.Buffer
	armorWriter, err := armor.Encode(&tokenBuffer, "PGP MESSAGE", map[string]string{})
	if err != nil {
		return
	}

	w, err := openpgp.Encrypt(armorWriter, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return
	}

	w.Write([]byte(token))
	w.Close()

	armorWriter.Close()

	encrypted = tokenBuffer.String()
	return
}

func (api *Api) Auth(ctx *macaron.Context, req AuthReq) {
	if req.GrantType != GrantPassword {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{400},
			Error: "invalid_grant",
			ErrorDescription: "GrantType must be set to password",
		})
		return
	}

	user, err := api.backend.Auth(req.Username, req.Password)
	if err != nil {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{401},
			Error: "invalid_grant",
			ErrorDescription: err.Error(),
		})
		return
	}

	sessionToken := "access_token"

	encryptedToken, err := encrypt(user, sessionToken)
	if err != nil {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{500},
			Error: "invalid_key",
			ErrorDescription: err.Error(),
		})
		return
	}

	api.sessions[sessionToken] = user.ID

	ctx.JSON(200, &AuthResp{
		Resp: Resp{1000},
		AccessToken: encryptedToken,
		ExpiresIn: 360000,
		TokenType: TokenBearer,
		Scope: "full mail payments reset keys",
		Uid: "uid",
		RefreshToken: "refresh_token",
		PrivateKey: user.EncPrivateKey,
		EncPrivateKey: user.EncPrivateKey,
		EventID: "gnFPgsx4P9uXvB7IW8sIAUEcxEGGGH7mmRFiCmWwcn1jY3hxPxnCh39qvQInv5LkQFPn5rYh8qzfP_bJPrvHrg==",
	})
}

func (api *Api) AuthCookies(ctx *macaron.Context, req AuthCookiesReq) {
	sessionToken := api.getSessionToken(ctx)
	if sessionToken == "" {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{401},
			Error: "invalid_session",
			ErrorDescription: "Not logged in",
		})
		return
	}

	if req.GrantType != GrantRefreshToken {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{400},
			Error: "invalid_grant",
			ErrorDescription: "GrantType must be set to refresh_token",
		})
		return
	}

	authCookie, _ := json.Marshal(&AuthCookie{
		AccessToken: sessionToken,
		Uid: "uid",
	})
	ctx.SetCookie("AUTH-" + sessionToken, string(authCookie), 0, "/api/", "", false, true)

	ctx.JSON(200, &AuthCookiesResp{
		Resp: Resp{1000},
		SessionToken: sessionToken,
	})
}
