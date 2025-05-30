package tools

import (
	"reflect"

	"github.com/sashabaranov/go-openai"
	"github.com/swaggest/jsonschema-go"
)

func OpenAITool(tool Tool) (openai.Tool, error) {
	emptyRequest := reflect.New(tool.RequestType())
	reflector := jsonschema.Reflector{}
	schema, err := reflector.Reflect(emptyRequest.Interface())
	if err != nil {
		return openai.Tool{}, err
	}

	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  schema,
		},
	}, nil
}
