package payment

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

var statusMap = map[stripe.PaymentIntentStatus]shop.OrderStatus{
	stripe.PaymentIntentStatusCanceled:              shop.Cancelled,
	stripe.PaymentIntentStatusProcessing:            shop.Pending,
	stripe.PaymentIntentStatusRequiresAction:        shop.Pending,
	stripe.PaymentIntentStatusRequiresCapture:       shop.Pending,
	stripe.PaymentIntentStatusRequiresConfirmation:  shop.Pending,
	stripe.PaymentIntentStatusRequiresPaymentMethod: shop.Pending,
	stripe.PaymentIntentStatusSucceeded:             shop.Paid,
}

type StripeClient struct {}

func NewStripeClient() *StripeClient {
	stripe.Key = viper.GetString("STRIPE_API_KEY")
	return &StripeClient{}
}

func (c *StripeClient) CreatePaymentIntentForOrder(ctx context.Context, order *shop.Order) error {
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
		log.Default().Errorf("stripe intent creation failed for order %s: %v", pi.ID, err)
		return err
	}

	order.PaymentIntentID = &pi.ID
	order.Status = statusMap[pi.Status]
	order.ClientSecret = &pi.ClientSecret

	return nil
}
