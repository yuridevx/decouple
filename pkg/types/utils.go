package types

import (
	"context"
	"reflect"
)

var contextType = reflect.TypeOf(new(context.Context)).Elem()
var errType = reflect.TypeOf(new(error)).Elem()

func IsContext(t reflect.Type) bool {
	if t.Kind() != reflect.Interface {
		return false
	}
	if t.Implements(contextType) {
		return true
	}
	return false
}

func UnwrapError(v reflect.Value) error {
	intf := v.Interface()
	if intf != nil {
		return intf.(error)
	}
	return nil
}

func IsError(t reflect.Type) bool {
	if t.Kind() != reflect.Interface {
		return false
	}
	if t.Implements(errType) {
		return true
	}
	return false
}
