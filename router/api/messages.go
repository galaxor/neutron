package api

import (
	"gopkg.in/macaron.v1"
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
	Response Resp
}

func (api *Api) GetMessagesCount(ctx *macaron.Context) {
	api.GetConversationsCount(ctx) // TODO?
}

func (api *Api) SetMessagesRead(ctx *macaron.Context, req SetMessagesReadReq) {
	userId := api.getUserId(ctx)
	resps := []*SetMessageReadResp{}

	for _, id := range req.IDs {
		// TODO
	}

	ctx.JSON(200, &SetMessagesReadResp{
		Resp: Resp{1001},
		Responses: resps,
	})
}
