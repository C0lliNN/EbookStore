package catalog

import "time"

type Book struct {
	ID               string
	Title            string
	Description      string
	AuthorName       string
	Images           []Image
	ContentBucketKey string
	Price            int
	ReleaseDate      time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Image struct {
	ID          string
	Description string
	BookID      string
}
