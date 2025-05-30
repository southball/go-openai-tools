package tools

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/sashabaranov/go-openai"
)

// Handles a single [openai.ToolCall] using the provided [Tool].
func HandleToolCall(tool Tool, toolCall openai.ToolCall) (openai.ChatCompletionMessage, error) {
	emptyRequest := reflect.New(tool.RequestType())
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), emptyRequest.Interface())
	if err != nil {
		return openai.ChatCompletionMessage{}, fmt.Errorf("failed to unmarshal tool call arguments: %w", err)
	}

	result, err := tool.CallFunction(emptyRequest.Elem().Interface())
	if err != nil {
		return openai.ChatCompletionMessage{}, fmt.Errorf("failed to call tool function: %w", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return openai.ChatCompletionMessage{}, fmt.Errorf("failed to marshal tool function call result: %w", err)
	}

	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		ToolCallID: toolCall.ID,
		Name:       tool.Name(),
		Content:    string(resultJSON),
	}, nil
}
