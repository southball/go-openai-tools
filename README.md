# go-openai-tools

[![Go Reference](https://pkg.go.dev/badge/southball/go-openai-tools.svg)](https://pkg.go.dev/southball/go-openai-tools) [![CI](https://github.com/southball/go-openai-tools/actions/workflows/ci.yaml/badge.svg)](https://github.com/southball/go-openai-tools/actions/workflows/ci.yaml)

This is a library to make tool definition and tool calls easy. The library is intended to be used with [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai).

## Installation

```sh
go get -u github.com/southball/go-openai-tools
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	openai "github.com/sashabaranov/go-openai"
	tools "github.com/southball/go-openai-tools"
)

type GetWeatherRequest struct {
	Location string `json:"location" required:"true" description:"City and country e.g. Bogotá, Colombia"`
}

type GetWeatherResponse struct {
	Weather string `json:"weather"`
}

func GetWeather(_ GetWeatherRequest) (GetWeatherResponse, error) {
	return GetWeatherResponse{Weather: "sunny"}, nil
}

const MODEL = "gpt-4o"

func main() {
	toolset, err := tools.NewToolSet(
		tools.F("GetWeather", "Get the current weather for a location.", GetWeather),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create toolset: %w", err))
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	initialMessages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful assistant."},
		{Role: openai.ChatMessageRoleUser, Content: "What is the weather like in Bogotá, Colombia?"},
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    MODEL,
			Messages: initialMessages,
			Tools:    toolset.OpenAITools(),
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create chat completion: %w", err))
	}

	answer := resp.Choices[0].Message
	toolCallMessages, err := toolset.HandleToolCalls(context.Background(), answer.ToolCalls)
	if err != nil {
		panic(fmt.Errorf("failed to handle tool calls: %w", err))
	}

	resp, err = client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    MODEL,
			Messages: slices.Concat(initialMessages, []openai.ChatCompletionMessage{answer}, toolCallMessages),
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create chat completion: %w", err))
	}

	fmt.Printf("Final response: %s\n", resp.Choices[0].Message.Content)
}
```
