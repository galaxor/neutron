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

func (api *Api) GetUserOrganization(ctx *macaron.Context) (err error) {
	userId := api.getUserId(ctx)

	user, err := api.backend.GetUser(userId)
	if err != nil {
		return err
	}

	domains, err := api.backend.ListDomains()
	if err != nil {
		return err
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
	return
}
