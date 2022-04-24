//go:build unit
// +build unit

package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilenameGenerator_NewUniqueName(t *testing.T) {
	name := FilenameGenerator(0).NewUniqueName("clean_code")

	assert.Regexp(t, ".+_\\d", name)
}
