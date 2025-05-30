package tools

import "reflect"

type typedTool[U, V any] interface {
	// Name of the tool.
	Name() string
	// Description of the tool.
	Description() string
	// Accepts argument of the type U and returns a value of type V or an error.
	CallFunction(args U) (V, error)
}

type Tool interface {
	// Name of the tool.
	Name() string
	// Description of the tool.
	Description() string
	// Type of the request.
	// The type should have struct tags as specified in https://pkg.go.dev/github.com/swaggest/jsonschema-go#Reflector.Reflect.
	RequestType() reflect.Type
	// Accepts argument of the type returned by RequestType(). Returns any value or an error.
	CallFunction(args any) (any, error)
}
