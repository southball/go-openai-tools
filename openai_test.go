package tools_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sashabaranov/go-openai"
	tools "github.com/southball/go-openai-tools"
)

func TestOpenAITool(t *testing.T) {
	testTool := tools.F("DemoTool", "A demo tool for OpenAI integration.", DemoTool)

	openAITestTool, err := tools.OpenAITool(testTool)
	if err != nil {
		t.Fatalf("Failed to convert tool to OpenAI format: %v", err)
	}

	expected := openai.Tool{
		Type: "function",
		Function: &openai.FunctionDefinition{
			Name:        "DemoTool",
			Description: "A demo tool for OpenAI integration.",
			Parameters:  json.RawMessage(`{"required":["arg1"],"properties":{"arg1":{"description":"First argument","type":"string"},"arg2":{"description":"Second argument","type":"integer"}},"type":"object"}`),
		},
	}
	if diff := cmp.Diff(expected, openAITestTool); diff != "" {
		t.Errorf("OpenAI tool does not match expected (-want +got):\n%s", diff)
	}
}
