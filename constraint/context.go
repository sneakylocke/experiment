package constraint

import (
	"github.com/juju/errors"
)

// Context is an interface to create objects that act as dictionaries. The Resolver will use a Context to pull
// information to Resolve whether or not a constraint has been satisfied.
type Context interface {
	value(key string) (interface{}, error)
}

// MapContext is a dictionary backed implementation of a Context
type MapContext struct {
	context map[string]interface{}
}

// NewMapContext returns a pointer to a MapContext
func NewMapContext(context map[string]interface{}) *MapContext {
	c := MapContext{}
	c.context = context
	return &c
}

func (context *MapContext) value(key string) (interface{}, error) {
	if value, ok := context.context[key]; ok {
		return value, nil
	}

	return nil, errors.Errorf("Key '%s' does not exists in context", key)
}
