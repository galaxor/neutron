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

	for _, kp := range req.Keys {
		err := api.backend.UpdateKeypair(userId, req.Password, kp)
		if err != nil {
			ctx.JSON(500, newErrorResp(err))
			return
		}
	}

	ctx.JSON(200, Resp{Ok})
}
