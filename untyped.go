package tools

import (
	"fmt"
	"reflect"
)

type untypedToolAdapter[U, V any] struct {
	tool typedTool[U, V]
}

func (u *untypedToolAdapter[U, V]) Name() string        { return u.tool.Name() }
func (u *untypedToolAdapter[U, V]) Description() string { return u.tool.Description() }
func (u *untypedToolAdapter[U, V]) RequestType() reflect.Type {
	return reflect.TypeOf(u.tool.CallFunction).In(0)
}
func (u *untypedToolAdapter[U, V]) CallFunction(args any) (any, error) {
	var expected U
	typedArgs, ok := reflect.ValueOf(args).Interface().(U)
	if !ok {
		return nil, fmt.Errorf("invalid argument type: expected type %T, got %T", expected, args)
	}
	result, err := u.tool.CallFunction(typedArgs)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Scrubs the type information from the tool
func t[U, V any](tool typedTool[U, V]) Tool {
	return &untypedToolAdapter[U, V]{tool: tool}
}
