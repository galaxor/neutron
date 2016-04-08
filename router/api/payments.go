package api

import (
	"gopkg.in/macaron.v1"
)

type PlansResp struct {
	Resp
	Plans []*Plan
}

type Plan struct {
	ID string
	Type int
	Cycle int
	Name string
	Currency string
	Amount int
	MaxDomains int
	MaxAddresses int
	MaxSpace int
	MaxMembers int
	TwoFactor int
}

func (api *Api) GetPlans(ctx *macaron.Context) {
	ctx.JSON(200, &PlansResp{
		Resp: Resp{Ok},
		Plans: []*Plan{},
	})
}

func (api *Api) GetSubscription(ctx *macaron.Context) {
	ctx.JSON(200, &ErrorResp{
		Resp: Resp{22110},
		Error: "You do not have an active subscription",
	})
}

type PaymentMethodsResp struct {
	Resp
	PaymentMethods []interface{} // TODO
}

func (api *Api) GetPaymentMethods(ctx *macaron.Context) {
	ctx.JSON(200, &PaymentMethodsResp{
		Resp: Resp{Ok},
		PaymentMethods: []interface{}{},
	})
}