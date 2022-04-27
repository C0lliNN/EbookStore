//go:build unit
// +build unit

package shop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrder_Complete(t *testing.T) {
	order := Order{}

	order.Complete()

	assert.Equal(t, Paid, order.Status)
}
