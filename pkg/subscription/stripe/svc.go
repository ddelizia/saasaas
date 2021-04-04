package stripe

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ddelizia/saasaas/pkg/ctx"
	"github.com/ddelizia/saasaas/pkg/subscription"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/stripe/stripe-go/v72"
	billingportalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/webhook"
)

type Service interface {
	subscription.Service
	CreateCheckoutSession(c context.Context, priceID string) (*stripe.CheckoutSession, error)
	CheckoutSession(c context.Context, sessionID string) (*stripe.CheckoutSession, error)
	CustomerPortal(c context.Context, sessionID string) (*stripe.BillingPortalSession, error)
	Webhook(c context.Context, r *http.Request)
	SaveCustomer(c context.Context, customerID string) error
	GetCustomer(c context.Context) (Customer, error)
}

type service struct {
}

func (s *service) AvailablePlans(c context.Context) ([]*subscription.Plan, error) {
	params := &stripe.ProductListParams{

	}
	iter := product.List(params)
	
	iter.
	
}

func (s *service) AccountPlan(c context.Context) (*Plan, error) {
	account, err := ctx.GetFromContext(c, ctx.AccountIDContextField)
	if err != nil {

	}

	
}

func (s *service) CreateCheckoutSession(c context.Context, priceID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: t.NewString("https://example.com/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  t.NewString("https://example.com/canceled.html"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price: stripe.String(priceID),
				// For metered billing, do not pass quantity
				Quantity: stripe.Int64(1),
			},
		},
	}

	session, err := session.New(params)
	if err != nil {

	}
	return session, nil
}

func (s *service) CheckoutSession(c context.Context, sessionID string) (*stripe.CheckoutSession, error) {
	result, err := session.Get(sessionID, nil)
	return result, err
}

func (s *service) CustomerPortal(c context.Context, sessionID string) (*stripe.BillingPortalSession, error) {
	ses, err := session.Get(sessionID, nil)
	if err != nil {

	}

	returnURL := os.Getenv("DOMAIN")

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(ses.Customer.ID),
		ReturnURL: stripe.String(returnURL),
	}
	ps, err := billingportalsession.New(params)
	if err != nil {

	}
	return ps, nil

}

func (s *service) Webhook(c context.Context, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {

		return
	}

	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {

		return
	}

	switch event.Type {
	case "checkout.session.completed":
		// Payment is successful and the subscription is created.
		// You should provision the subscription.
	case "invoice.paid":
		// Continue to provision the subscription as payments continue to be made.
		// Store the status in your database and check when a user accesses your service.
		// This approach helps you avoid hitting rate limits.
	case "invoice.payment_failed":
		// The payment failed or the customer does not have a valid payment method.
		// The subscription becomes past_due. Notify your customer and send them to the
		// customer portal to update their payment information.
	default:
		// unhandled event type
	}
}
