package catalog

import "time"

type Book struct {
	ID                   string
	Title                string
	Description          string
	AuthorName           string
	PosterImageBucketKey string
	PosterImageLink      string `gorm:"-"`
	ContentBucketKey     string
	Price                int
	ReleaseDate          time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (b *Book) SetPosterImageLink(imageLink string) {
	b.PosterImageLink = imageLink
}
