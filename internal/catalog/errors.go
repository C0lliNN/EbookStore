package catalog

import "fmt"

var ErrForbiddenCatalogAccess = fmt.Errorf("the access to this action is restricted to allowed users")
