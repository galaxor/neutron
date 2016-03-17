package api

import (
	"gopkg.in/macaron.v1"

	"github.com/emersion/neutron/backend"
)

type SetMessagesReadReq struct {
	Req
	IDs []string
}

type SetMessagesReadResp struct {
	Resp
	Responses []*SetMessageReadResp
}

type SetMessageReadResp struct {
	ID string
	Response interface{}
}

func (api *Api) GetMessagesCount(ctx *macaron.Context) {
	api.GetConversationsCount(ctx) // TODO?
}

func (api *Api) setMessagesRead(ctx *macaron.Context, ids []string, value int) {
	userId := api.getUserId(ctx)

	var resps []*SetMessageReadResp

	for _, id := range ids {
		err := api.backend.UpdateMessage(userId, &backend.MessageUpdate{
			Message: &backend.Message{
				ID: id,
				IsRead: 1,
			},
			IsRead: true,
		})

		r := &SetMessageReadResp{ ID: id }
		if err != nil {
			r.Response = &ErrorResp{
				Resp: Resp{500},
				ErrorDescription: err.Error(),
			}
		} else {
			r.Response = &Resp{1000}
		}
		resps = append(resps, r)
	}

	ctx.JSON(200, &SetMessagesReadResp{
		Resp: Resp{1001},
		Responses: resps,
	})
}

func (api *Api) SetMessagesRead(ctx *macaron.Context, req SetMessagesReadReq) {
	api.setMessagesRead(ctx, req.IDs, 1)
}

func (api *Api) SetMessagesUnread(ctx *macaron.Context, req SetMessagesReadReq) {
	api.setMessagesRead(ctx, req.IDs, 0)
}
