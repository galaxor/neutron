package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type UserResp struct {
	Resp
	User *backend.User
}

func GetCurrentUser(ctx *macaron.Context) {
	user, err := backend.Get(userId)
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
