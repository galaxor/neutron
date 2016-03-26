package api

import (
	"gopkg.in/macaron.v1"
)

func (api *Api) GetAttachment(ctx *macaron.Context) (b []byte, err error) {
	userId := api.getUserId(ctx)
	id := ctx.Params("id")

	att, b, err := api.backend.ReadAttachment(userId, id)
	if err != nil {
		return
	}

	ctx.Resp.Header().Set("Content-Type", att.MIMEType)
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename=\""+att.Name+"\"")
	ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Resp.Header().Set("Expires", "0")
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.Header().Set("Pragma", "public")
	return
}
