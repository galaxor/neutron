package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

func (api *Api) GetMessagesCount(ctx *macaron.Context) {
	api.GetConversationsCount(ctx) // TODO?
}

func (api *Api) setMessagesRead(ctx *macaron.Context, req BatchReq, value int) {
	userId := api.getUserId(ctx)

	var respItems []*BatchRespItem

	for _, id := range req.IDs {
		err := api.backend.UpdateMessage(userId, &backend.MessageUpdate{
			Message: &backend.Message{
				ID: id,
				IsRead: 1,
			},
			IsRead: true,
		})

		r := &BatchRespItem{ ID: id }
		if err != nil {
			r.Response = &ErrorResp{
				Resp: Resp{InternalServerError},
				ErrorDescription: err.Error(),
			}
		} else {
			r.Response = &Resp{Ok}
		}
		respItems = append(respItems, r)
	}

	ctx.JSON(200, &BatchResp{
		Resp: Resp{Batch},
		Responses: respItems,
	})
}

func (api *Api) SetMessagesRead(ctx *macaron.Context, req BatchReq) {
	api.setMessagesRead(ctx, req, 1)
}

func (api *Api) SetMessagesUnread(ctx *macaron.Context, req BatchReq) {
	api.setMessagesRead(ctx, req, 0)
}
