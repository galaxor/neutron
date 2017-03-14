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

type SubscriptionResp struct {
	Resp
	Subscription *Subscription
}

type Subscription struct {
	ID string
	InvoiceID string
	Cycle int
	PeriodStart int64
	PeriodEnd int64
	CouponCode string
	Currency string
	Amount int
	Plans []*Plan
}

func (api *Api) GetPlans(ctx *macaron.Context) {
	ctx.JSON(200, &PlansResp{
		Resp: Resp{Ok},
		Plans: []*Plan{},
	})
}

func (api *Api) GetSubscription(ctx *macaron.Context) {
	ctx.JSON(200, &SubscriptionResp{
		Resp: Resp{Ok},
		Subscription: &Subscription{
			ID: "0",
			InvoiceID: "0",
			Plans: []*Plan{},
		},
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
