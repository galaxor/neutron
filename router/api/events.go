package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type EventResp struct {
	Resp
	*backend.Event
}

func (api *Api) GetEvent(ctx *macaron.Context) (err error) {
	userId := api.getUserId(ctx)
	eventId := ctx.Params("event")

	event, err := api.backend.GetEventsAfter(userId, eventId)
	if err != nil {
		return
	}

	// Client crashes if Notices is null
	if event.Notices == nil {
		event.Notices = []string{}
	}

	ctx.JSON(200, &EventResp{
		Resp: Resp{Ok},
		Event: event,
	})
	return
}
