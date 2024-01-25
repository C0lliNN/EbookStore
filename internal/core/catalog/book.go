package catalog

import "time"

type Book struct {
	ID          string
	Title       string
	Description string
	AuthorName  string
	Images      []Image
	ContentID   string
	Price       int
	ReleaseDate time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Image struct {
	ID          string
	Description string
	BookID      string
}

func (b Book) MainImageID() string {
	// TODO: implement position to be able to choose main image
	if len(b.Images) > 0 {
		return b.Images[0].ID
	}
	return ""
}
