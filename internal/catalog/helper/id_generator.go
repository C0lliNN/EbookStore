package helper

import "github.com/google/uuid"

type IDGenerator byte

func NewIDGenerator() IDGenerator {
	return IDGenerator(0)
}

func (g IDGenerator) NewID() string {
	return uuid.NewString()
}
