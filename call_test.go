package tools_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sashabaranov/go-openai"
	tools "github.com/southball/go-openai-tools"
)

func TestHandleToolCall(t *testing.T) {
	testTool := tools.F("DemoTool", "A demo tool for OpenAI integration.", DemoTool)

	testToolCall := openai.ToolCall{
		ID:    "test-tool-call-id",
		Type:  openai.ToolTypeFunction,
		Index: ptr(0),
		Function: openai.FunctionCall{
			Name:      "DemoTool",
			Arguments: `{"arg1": "test", "arg2": 42}`,
		},
	}

	handledMessage, err := tools.HandleToolCall(testTool, testToolCall)
	if err != nil {
		t.Fatalf("Failed to handle tool call: %v", err)
	}

	expectedMessage := openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		ToolCallID: "test-tool-call-id",
		Name:       "DemoTool",
		Content:    `{"result":"Processed: test and 42"}`,
	}

	if diff := cmp.Diff(expectedMessage, handledMessage); diff != "" {
		t.Errorf("Handled message does not match expected (-want +got):\n%s", diff)
	}
}
