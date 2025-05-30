package tools

import (
	"encoding/json"
	"reflect"

	"github.com/sashabaranov/go-openai"
	"github.com/swaggest/jsonschema-go"
)

// Convert a [Tool] to an [openai.Tool].
func OpenAITool(tool Tool) (openai.Tool, error) {
	emptyRequest := reflect.New(tool.RequestType())
	reflector := jsonschema.Reflector{}
	schema, err := reflector.Reflect(emptyRequest.Interface())
	if err != nil {
		return openai.Tool{}, err
	}

	schemaJson, err := json.Marshal(schema)
	if err != nil {
		return openai.Tool{}, err
	}

	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  json.RawMessage(schemaJson),
		},
	}, nil
}
