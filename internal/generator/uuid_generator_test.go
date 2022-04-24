//go:build unit
// +build unit

package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUIDGenerator_NewID(t *testing.T) {
	id := UUIDGenerator{}.NewID()

	assert.Len(t, id, 36)
}
