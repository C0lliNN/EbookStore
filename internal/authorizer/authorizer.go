package authorizer

import (
	"context"
	"fmt"
	"github.com/casbin/casbin"
)

var ErrUnauthorized = fmt.Errorf("the current action is not allowed by the current user")

type Config struct {
	// Model represents a filepath
	Model string
	// Policy represents a filepath
	Policy string
}

type Authorizer struct {
	*casbin.Enforcer
}

func New(c Config) *Authorizer {
	return &Authorizer{
		Enforcer: casbin.NewEnforcer(c.Model, c.Policy),
	}
}

func (a *Authorizer) Authorize(ctx context.Context, object interface{}, action string) error {
	if !a.Enforce(subject, object, action) {
		return ErrUnauthorized
	}

	return nil
}
