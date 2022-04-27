package shop

import "fmt"

var ErrOrderNotPaid = fmt.Errorf("only books from paid orders can be downloaded")
