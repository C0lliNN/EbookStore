package helper

import (
	"fmt"
	"time"
)

type FilenameGenerator byte

func NewFilenameGenerator() FilenameGenerator {
	return 0
}

func (g FilenameGenerator) NewUniqueName(filename string) string {
	unix := time.Now().UTC().Unix()

	return fmt.Sprintf("%s_%d", filename, unix)
}
