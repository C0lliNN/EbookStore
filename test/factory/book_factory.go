package factory

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"time"
)

func NewBook() catalog.Book {
	return catalog.Book{
		ID:                   faker.UUIDHyphenated(),
		Title:                faker.TitleMale(),
		Description:          faker.Sentence(),
		AuthorName:           faker.FirstName(),
		PosterImageBucketKey: faker.UUIDHyphenated(),
		ContentBucketKey:     faker.UUIDHyphenated(),
		Price:                40000,
		ReleaseDate:          time.Unix(faker.RandomUnixTime(), 0).UTC(),
		CreatedAt:            time.Unix(faker.RandomUnixTime(), 0).UTC(),
		UpdatedAt:            time.Unix(faker.RandomUnixTime(), 0).UTC(),
	}
}

func NewPaginatedBooks() catalog.PaginatedBooks {
	return catalog.PaginatedBooks{
		Books:      []catalog.Book{NewBook(), NewBook()},
		Limit:      10,
		Offset:     0,
		TotalBooks: 2,
	}
}
