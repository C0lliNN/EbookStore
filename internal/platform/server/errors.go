package server

type BindingErr struct {
	Err error
}

func (e *BindingErr) Error() string {
	return e.Err.Error()
}
