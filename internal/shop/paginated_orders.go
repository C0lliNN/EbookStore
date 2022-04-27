package shop

type PaginatedOrders struct {
	Orders      []Order
	Limit       int
	Offset      int
	TotalOrders int64
}
