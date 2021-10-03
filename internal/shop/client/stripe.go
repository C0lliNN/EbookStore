package client

import (
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

var statusMap = map[stripe.PaymentIntentStatus]model.OrderStatus{
	stripe.PaymentIntentStatusCanceled:              model.Cancelled,
	stripe.PaymentIntentStatusProcessing:            model.Pending,
	stripe.PaymentIntentStatusRequiresAction:        model.Pending,
	stripe.PaymentIntentStatusRequiresCapture:       model.Pending,
	stripe.PaymentIntentStatusRequiresConfirmation:  model.Pending,
	stripe.PaymentIntentStatusRequiresPaymentMethod: model.Pending,
	stripe.PaymentIntentStatusSucceeded:             model.Pending,
}

type StripeClient byte

func NewStripeClient() StripeClient {
	stripe.Key = viper.GetString("STRIPE_API_KEY")
	return StripeClient(0)
}

func (c StripeClient) CreatePaymentIntentForOrder(order *model.Order) error {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(order.Total),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}

	params.AddMetadata("orderID", order.ID)
	params.AddMetadata("userID", order.UserID)
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Logger.Errorf("stripe intent creation failed for order %s: %v", pi.ID, err)
		return err
	}

	order.PaymentIntent = &pi.ID
	order.Status = statusMap[pi.Status]
	order.ClientSecret = &pi.ClientSecret

	return nil
}
