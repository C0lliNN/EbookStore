package query

// Page is a struct that represents a page in a query independent of the database.
type Page struct {
	Number int
	Size int
}

// DefaultPage configuration
var DefaultPage = Page {
	Number: 1,
	Size: 15,
}

// Offset returns the offset of the page.
func (p Page) Offset() int {
	return (p.Number - 1) * p.Size
}
