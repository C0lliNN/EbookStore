package auth

import "fmt"

var ErrWrongPassword = fmt.Errorf("the provided password is incorrect")
