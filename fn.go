package tools

type fnTool[U, V any] struct {
	name        string
	description string
	callFunc    func(args U) (V, error)
}

func (f *fnTool[U, V]) Name() string                   { return f.name }
func (f *fnTool[U, V]) Description() string            { return f.description }
func (f *fnTool[U, V]) CallFunction(args U) (V, error) { return f.callFunc(args) }

// Create a [Tool] with the given name, description, and function.
func F[U, V any](name string, description string, callFunc func(args U) (V, error)) Tool {
	return t(&fnTool[U, V]{name, description, callFunc})
}
