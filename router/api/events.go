package api

import (
	"gopkg.in/macaron.v1"
)

type EventResp struct {
	Resp
	EventID string
	Refresh int
	Reload int
	Notices []string

	// See https://github.com/ProtonMail/WebClient/blob/118b9473c837eaa1fb4dc9e2591013b24bedfcbe/src/app/services/event.js#L274
	//Labels
	//Contacts
	//User
	//Messages
	//Conversations
	//MessageCounts
	//ConversationCounts
	//UsedSpace
	//Domains
	//Members
	//Organization
}

func (api *Api) GetEvent(ctx *macaron.Context) {
	id := ctx.Params("event")

	ctx.JSON(200, &EventResp{
		Resp: Resp{1000},
		EventID: id,
		Notices: []string{},
	})
}
