package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBook_MainImageID(t *testing.T) {
	tests := []struct {
		Name     string
		Book     Book
		Expected string
	}{
		{
			Name:     "no images",
			Book:     Book{},
			Expected: "",
		},
		{
			Name: "one image",
			Book: Book{
				Images: []Image{
					{ID: "image-id"},
				},
			},
			Expected: "image-id",
		},
		{
			Name: "two images",
			Book: Book{
				Images: []Image{
					{ID: "image-id"},
					{ID: "image-id2"},
				},
			},
			Expected: "image-id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actual := tc.Book.MainImageID()
			assert.Equal(t, tc.Expected, actual)
		})
	}
}
