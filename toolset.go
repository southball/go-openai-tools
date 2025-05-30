package tools

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/sashabaranov/go-openai"
)

// A collection of [Tool].
type ToolSet struct {
	tools       map[string]Tool
	openaiTools []openai.Tool
}

// Construct a [ToolSet] from the given list of [Tool].
func NewToolSet(toolset ...Tool) (*ToolSet, error) {
	tools := make(map[string]Tool, len(toolset))
	for _, tool := range toolset {
		if _, exists := tools[tool.Name()]; exists {
			return nil, fmt.Errorf("duplicate tool name: %s", tool.Name())
		}
		tools[tool.Name()] = tool
	}

	openAITools := make([]openai.Tool, 0, len(tools))
	for _, tool := range tools {
		openAITool, err := OpenAITool(tool)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tool %s to OpenAI tool: %w", tool.Name(), err)
		}
		openAITools = append(openAITools, openAITool)
	}

	return &ToolSet{tools, openAITools}, nil
}

// Process a list of [openai.ToolCall] and return a list of [openai.ChatCompletionMessage] corresponding to the results of those calls.
func (t ToolSet) HandleToolCalls(
	ctx context.Context,
	toolCalls []openai.ToolCall,
) ([]openai.ChatCompletionMessage, error) {
	messageCh := make(chan openai.ChatCompletionMessage, len(toolCalls))
	errCh := make(chan error, len(toolCalls))

	var wg sync.WaitGroup
	wg.Add(len(toolCalls))
	for _, toolCall := range toolCalls {
		go func() {
			defer wg.Done()

			tool, exists := t.tools[toolCall.Function.Name]
			if !exists {
				errCh <- fmt.Errorf("tool not found: %s", toolCall.Function.Name)
				return
			}

			message, err := HandleToolCall(tool, toolCall)
			if err != nil {
				errCh <- err
				return
			}

			messageCh <- message
		}()
	}

	go func() {
		wg.Wait()
		close(messageCh)
		close(errCh)
	}()

	messages := make([]openai.ChatCompletionMessage, 0, len(toolCalls))
	errs := make([]error, 0)
	for message := range messageCh {
		messages = append(messages, message)
	}
	for err := range errCh {
		errs = append(errs, err)
	}

	return messages, errors.Join(errs...)
}

// Get a list of [openai.Tool] from the [ToolSet].
func (t ToolSet) OpenAITools() []openai.Tool {
	return t.openaiTools
}
