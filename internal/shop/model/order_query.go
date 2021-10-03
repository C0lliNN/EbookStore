package model

type OrderQuery struct {
	Status OrderStatus
	BookID string
	UserID string
	Limit  int
	Offset int
}

func (q OrderQuery) CreateCriteria() []Criteria {
	criteria := make([]Criteria, 0)

	criteria = append(criteria, NewEqualCriteria("status", string(q.Status)))
	criteria = append(criteria, NewEqualCriteria("book_id", q.BookID))
	criteria = append(criteria, NewEqualCriteria("user_id", q.UserID))

	return criteria
}
