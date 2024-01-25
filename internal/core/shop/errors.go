package shop

import "fmt"

var ErrOrderNotCompleted = fmt.Errorf("only books from completed orders can be downloaded")
var ErrForbiddenOrderAccess = fmt.Errorf("the access to this order is restricted to allowed users")
var ErrItemAlreadyInCart = fmt.Errorf("item already in cart")
var ErrItemNotFoundInCart = fmt.Errorf("item not found in cart")
var ErrItemNotFoundInOrder = fmt.Errorf("item not found in order")
