package shop

type SearchOrders struct {
	Status  string `form:"status"`
	Page    int    `form:"page"`
	PerPage int    `form:"perPage"`
}

func (s *SearchOrders) OrderQuery() OrderQuery {
	if s.Page == 0 {
		s.Page = 1
	}

	if s.PerPage == 0 {
		s.PerPage = 10
	}

	return OrderQuery{
		Status: OrderStatus(s.Status),
		Limit:  s.PerPage,
		Offset: (s.Page - 1) * s.PerPage,
	}
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
