package api

import (
	"encoding/base64"

	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type UserResp struct {
	Resp
	User *backend.User
}

type CreateUserReq struct {
	Req
	Username string
	Password string
	Domain string
	News bool
	PrivateKey string
	Token string
	TokenType string
}

type DirectUserResp struct {
	Resp
	Direct int
}

type UsernameAvailableResp struct {
	Resp
	Available int
}

type UpdateUserDisplayNameReq struct {
	Req
	DisplayName string
}

func (api *Api) GetCurrentUser(ctx *macaron.Context) {
	userId := api.getUserId(ctx)

	user, err := api.backend.GetUser(userId)
	if err != nil {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{NotFound},
			Error: "invalid_user",
			ErrorDescription: err.Error(),
		})
		return
	}

	ctx.JSON(200, &UserResp{
		Resp: Resp{Ok},
		User: user,
	})
}

func (api *Api) CreateUser(ctx *macaron.Context, req CreateUserReq) (err error) {
	// TODO: support req.Domain, req.Token, req.TokenType

	user, err := api.backend.InsertUser(&backend.User{
		Name: req.Username,
		EncPrivateKey: req.PrivateKey,
	}, req.Password)
	if err != nil {
		return
	}

	ctx.JSON(200, &UserResp{
		Resp: Resp{Ok},
		User: user,
	})
	return
}

func (api *Api) GetDirectUser(ctx *macaron.Context) {
	ctx.JSON(200, &DirectUserResp{
		Resp: Resp{Ok},
		Direct: 1,
	})
}

func (api *Api) GetUsernameAvailable(ctx *macaron.Context) (err error) {
	username := ctx.Params("username")

	available, err := api.backend.IsUsernameAvailable(username)
	if err != nil {
		return
	}

	value := 0
	if available {
		value = 1
	}

	ctx.JSON(200, &UsernameAvailableResp{
		Resp: Resp{Ok},
		Available: value,
	})
	return
}

func (api *Api) UpdateUserDisplayName(ctx *macaron.Context, req UpdateUserDisplayNameReq) (err error) {
	err = api.backend.UpdateUser(&backend.UserUpdate{
		User: &backend.User{
			DisplayName: req.DisplayName,
		},
		DisplayName: true,
	})
	if err != nil {
		return
	}

	ctx.JSON(200, &Resp{Ok})
	return
}

func (api *Api) GetPublicKey(ctx *macaron.Context) (err error) {
	b, err := base64.URLEncoding.DecodeString(ctx.Params("email"))
	if err != nil {
		return
	}

	email := string(b)

	key, err := api.backend.GetPublicKey(email)
	if err != nil {
		return
	}

	resp := map[string]interface{}{ "Code": Ok }
	if key != "" {
		resp[email] = key
	}
	ctx.JSON(200, resp)
	return
}
