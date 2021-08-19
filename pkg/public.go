package pkg

import (
	"reflect"
)

type Decouple interface {
	// Request
	// Function interfaces can be
	// func (ctx context.Context, request Type) (Type, error)
	// func (request Type) (Type, error)
	// func (request Type) Type
	Request(fn interface{})

	// Subscribe
	// Function interfaces can be
	// func (message Type)
	// func (message Type) error
	Subscribe(fn interface{})

	SetProtoName(obj interface{}, name string)
	GetProtoName(obj interface{}) string

	SetTypeName(t reflect.Type, name string)
	GetTypeName(t reflect.Type) string
}

type Call interface {
	Request(request interface{}, options ...CallOption) (interface{}, error)
	Broadcast(notification interface{}, options ...CallOption) ([]interface{}, error)
}
