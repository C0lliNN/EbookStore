package catalog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBook_SetPosterImageLink(t *testing.T) {
	book := Book{}
	link := "some-link"

	book.SetPosterImageLink(link)

	assert.Equal(t, book.PosterImageLink, link)
}
