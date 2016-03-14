package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type LabelsResp struct {
	Resp
	Labels []*backend.Label
}

func (api *Api) GetLabels(ctx *macaron.Context) (err error) {
	userId := api.getUserId(ctx)

	labels, err := api.backend.GetLabels(userId)
	if err != nil {
		return
	}

	ctx.JSON(200, &LabelsResp{
		Resp: Resp{1000},
		Labels: labels,
	})
	return
}
