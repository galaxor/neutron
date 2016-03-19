package api

import (
	"gopkg.in/macaron.v1"
)

type AvailableDomainsResp struct {
	Resp
	Domains []string
}

func (api *Api) GetAvailableDomains(ctx *macaron.Context) (err error) {
	domains, err := api.backend.ListDomains()
	if err != nil {
		return
	}

	domainNames := make([]string, len(domains))
	for i, d := range domains {
		domainNames[i] = d.Name
	}

	ctx.JSON(200, &AvailableDomainsResp{
		Resp: Resp{Ok},
		Domains: domainNames,
	})
	return
}
