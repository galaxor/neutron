package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type UpdateUserPrivateKeyReq struct {
	Password string
	Keys []*backend.Keypair
}

func (api *Api) UpdateUserPrivateKey(ctx *macaron.Context, req UpdateUserPrivateKeyReq) {
	userId := api.getUserId(ctx)

	user, err := api.backend.GetUser(userId)
	if err != nil {
		ctx.JSON(500, newErrorResp(err))
		return
	}

	// Check password
	user, err = api.backend.Auth(user.Name, req.Password)
	if err != nil {
		ctx.JSON(500, newErrorResp(err))
		return
	}

	for _, kp := range req.Keys {
		email := kp.ID

		_, err := api.backend.UpdateKeypair(email, req.Password, kp)
		if err != nil {
			ctx.JSON(500, newErrorResp(err))
			return
		}
	}

	ctx.JSON(200, Resp{Ok})
}
