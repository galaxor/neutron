package api

import (
	"gopkg.in/macaron.v1"
)

type Organization struct {
	ID string
	UsedDomains int
	MaxDomains int
	UsedAddresses int
	MaxAddresses int
}

type OrganizationResp struct {
	Resp
	Organization *Organization
}

func (api *Api) GetUserOrganization(ctx *macaron.Context) {
	userId := api.getUserId(ctx)

	user, err := api.backend.GetUser(userId)
	if err != nil {
		ctx.JSON(200, newErrorResp(err))
		return
	}

	domains, err := api.backend.ListDomains()
	if err != nil {
		ctx.JSON(200, newErrorResp(err))
		return
	}

	org := &Organization{
		ID: user.ID,
		UsedDomains: len(domains),
		UsedAddresses: len(user.Addresses),
	}

	ctx.JSON(200, &OrganizationResp{
		Resp: Resp{Ok},
		Organization: org,
	})
}
