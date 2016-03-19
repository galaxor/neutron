package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type UpdateUserSettingsReq struct {
	Req
	*backend.User
	Password string
}

func (api *Api) updateUserSettings(ctx *macaron.Context, update *backend.UserUpdate, updated *backend.User) {
	updated.ID = api.getUserId(ctx)

	update.User = updated

	err := api.backend.UpdateUser(update)
	if err != nil {
		ctx.JSON(500, newErrorResp(err))
		return
	}

	ctx.JSON(200, &Resp{Ok})
	return
}

func (api *Api) UpdateUserDisplayName(ctx *macaron.Context, req UpdateUserSettingsReq) {
	api.updateUserSettings(ctx, &backend.UserUpdate{DisplayName: true}, req.User)
}

func (api *Api) UpdateUserSignature(ctx *macaron.Context, req UpdateUserSettingsReq) {
	api.updateUserSettings(ctx, &backend.UserUpdate{Signature: true}, req.User)
}
