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
	Email string
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

func populateUser(user *backend.User) {
	if user.EncPrivateKey == "" || user.PublicKey == "" {
		keyring := user.GetMainAddress().Keys[0] // TODO: find a better way to get the keyring
		user.EncPrivateKey = keyring.PrivateKey
		user.PublicKey = keyring.PublicKey
	}

	for _, addr := range user.Addresses {
		if addr.DisplayName == "" {
			addr.DisplayName = user.DisplayName
		}
		if len(addr.Keys) > 0 {
			addr.HasKeys = 1
		}
	}

	user.Role = backend.RolePaidAdmin
	user.Subscribed = 1
	user.Private = 1
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

	populateUser(user)

	ctx.JSON(200, &UserResp{
		Resp: Resp{Ok},
		User: user,
	})
}

func (api *Api) CreateUser(ctx *macaron.Context, req CreateUserReq) (err error) {
	// TODO: check req.Token & req.TokenType

	domain, err := api.backend.GetDomainByName(req.Domain)
	if err != nil {
		return
	}

	email := req.Username + "@" + domain.Name

	// Insert user

	user, err := api.backend.InsertUser(&backend.User{
		Name: req.Username,
		NotificationEmail: req.Email,
		Addresses: []*backend.Address{
			&backend.Address{
				DomainID: domain.ID,
				Email: email,
				Send: 1,
				Receive: 1,
				Status: 1,
				Type: 1,
			},
		},
	}, req.Password)
	if err != nil {
		return
	}

	// Insert keypair

	keypair := backend.NewKeypair("", req.PrivateKey)
	keypair, err = api.backend.UpdateKeypair(email, req.Password, keypair)
	if err != nil {
		return
	}

	user.GetMainAddress().Keys = []*backend.Keypair{keypair}
	populateUser(user)

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
	resp[email] = key
	ctx.JSON(200, resp)
	return
}
