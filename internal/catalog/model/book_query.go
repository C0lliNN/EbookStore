package model

type BookQuery struct {
	Title       string
	Description string
	AuthorName  string
	Limit       int
	Offset      int
}

func (q BookQuery) CreateCriteria() []Criteria {
	return []Criteria{
		NewILikeCriteria("title", q.Title),
		NewILikeCriteria("description", q.Description),
		NewEqualCriteria("author_name", q.AuthorName),
	}
}
