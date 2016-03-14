package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type UserResp struct {
	Resp
	User *backend.User
}

func (api *Api) GetCurrentUser(ctx *macaron.Context) {
	user, err := api.backend.GetUser(userId)
	if err != nil {
		ctx.JSON(200, &ErrorResp{
			Resp: Resp{404},
			Error: "invalid_user",
			ErrorDescription: err.Error(),
		})
		return
	}

	ctx.JSON(200, &UserResp{
		Resp: Resp{1000},
		User: user,
	})
}
