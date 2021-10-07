// +build unit

package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIDGenerator_NewID(t *testing.T) {
	id := NewIDGenerator().NewID()

	assert.Len(t, id, 36)
}
