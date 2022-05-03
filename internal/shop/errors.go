package shop

import "fmt"

var ErrOrderNotPaid = fmt.Errorf("only books from paid orders can be downloaded")
var ErrForbiddenOrderAccess = fmt.Errorf("the access to this order is restricted to allowed users")
