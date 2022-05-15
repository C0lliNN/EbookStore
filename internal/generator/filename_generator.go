package generator

import (
	"fmt"
	"time"
)

type FilenameGenerator struct{}

func NewFilenameGenerator() *FilenameGenerator {
	return &FilenameGenerator{}
}

func (g *FilenameGenerator) NewUniqueName(name string) string {
	unix := time.Now().UTC().Unix()

	return fmt.Sprintf("%s_%d", name, unix)
}
