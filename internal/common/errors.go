package common

type ErrNotValid struct {
	Input string
	Err   error
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

type ErrEntityNotFound struct {
	Entity string
	Err    error
}

func (e *ErrEntityNotFound) Error() string {
	return e.Entity + " could not be found: " + e.Err.Error()
}

type ErrWrongPassword struct {
	Err error
}

func (e *ErrWrongPassword) Error() string {
	return e.Err.Error()
}

type ErrOrderNotPaid struct {
	Err error
}

func (e *ErrOrderNotPaid) Error() string {
	return e.Err.Error()
}
