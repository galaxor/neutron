package api

import (
	"gopkg.in/macaron.v1"
)

type AvailableDomainsResp struct {
	Resp
	Domains []string
}

func (api *Api) GetAvailableDomains(ctx *macaron.Context) {
	ctx.JSON(200, &AvailableDomainsResp{
		Resp: Resp{Ok},
		Domains: []string{"example.org"},
	})
	return
}
