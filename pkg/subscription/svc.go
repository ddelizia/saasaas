package subscription

import "context"

//https://stripe.com/docs/billing/subscriptions/checkout

type Service interface {
	AvailablePlans(c context.Context) ([]*Plan, error)
	AccountPlan(c context.Context) (*Plan, error)
}
