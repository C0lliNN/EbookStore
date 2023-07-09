package shop

import "github.com/ebookstore/internal/core/query"

type SearchOrders struct {
	Status  string `form:"status"`
	Page    int    `form:"page"`
	PerPage int    `form:"perPage"`
}

func (s *SearchOrders) CreateQuery() query.Query {
	q := query.New()

	if s.Status != "" {
		q.And(query.Condition{Field: "status", Operator: query.Equal, Value: s.Status})
	}

	return *q
}

func (s *SearchOrders) CreatePage() query.Page {
	p := query.DefaultPage

	if (s.Page > 0) {
		p.Number = s.Page
	}

	if (s.PerPage > 0) {
		p.Size = s.PerPage
	}

	return p
}

type CreateOrder struct {
	BookID string `json:"bookId" validate:"required,max=36"`
}

func (c CreateOrder) Order(orderId, userId string) Order {
	return Order{
		ID:     orderId,
		BookID: c.BookID,
		UserID: userId,
	}
}

type HandleStripeWebhook struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Request  map[string]interface{} `json:"request"`
	Livemode bool                   `json:"livemode"`
	Created  int                    `json:"created"`
}
