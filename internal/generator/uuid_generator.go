package generator

import "github.com/google/uuid"

type UUIDGenerator struct{}

func NewUUIDGenerator() UUIDGenerator {
	return UUIDGenerator{}
}

func (g UUIDGenerator) NewID() string {
	return uuid.NewString()
}
