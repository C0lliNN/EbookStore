package catalog

type PaginatedBooks struct {
	Books      []Book
	Limit      int
	Offset     int
	TotalBooks int64
}
