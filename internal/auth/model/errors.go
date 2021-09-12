package model

type ErrNotValid struct {
	Input string
	Err error
}

func (e *ErrNotValid) Error() string {
	return e.Input + " not valid: " + e.Err.Error()
}
