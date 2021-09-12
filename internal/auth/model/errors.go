package model

type ErrNotValid struct {
	Input string
	Err error
}

func (e *ErrNotValid) Error() string {
	return e.Input + " not valid: " + e.Err.Error()
}

type ErrDuplicateKey struct {
	Key string
	Err error
}

func (e *ErrDuplicateKey) Error() string {
	return e.Key + " violation: " + e.Err.Error()
}
