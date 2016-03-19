package api

import (
	"gopkg.in/macaron.v1"
)

type OrganizationResp struct {
	Resp
	Organization interface{} // TODO
}

func (api *Api) GetUserOrganization(ctx *macaron.Context) {
	ctx.JSON(200, &OrganizationResp{
		Resp: Resp{Ok},
	})
}
