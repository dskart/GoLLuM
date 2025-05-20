# 🤖 OpenAI

A clean and idiomatic Go client for the OpenAI API, providing strong typing, error handling, and configurable behavior for seamless integration with LLM-powered applications.

## Features

- **Type-Safe API**: Fully typed interfaces for OpenAI's API endpoints
- **Configurable Client**: Extensive configuration options to customize behavior
- **Retryable HTTP Client**: Built-in support for retries with exponential backoff
- **Function Calling Support**: First-class support for OpenAI's function/tool calling capabilities
- **Flexible Logging**: Supports standard library logger, zap, and other logging integrations

## Installation

```bash
go get github.com/dskart/gollum/openai
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dskart/gollum/openai"
)

func main() {
	openaiCfg := openai.Config{
		GptModel:  "gpt-4",
		OpenAiKey: "FOO",
	}

	llm, err := openai.New(openaiCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize openai: %w", err)
	}

	msgs := []openai.Message{
		{
			Role:    openai.SystemRoleType,
			Content: openai.StrPtr("You are an expert at Math"),
		},
		{
			Role:    openai.UserRoleType,
			Content: openai.StrPtr("What is 1+1"),
		},
	}
	opts := []func(*openai.ChatCompletionOptions){}
	resp, err := llm.ChatCompletionCreate(ctx, msgs, opts...)
	if err != nil {
		return fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("no choices returned")
	}

	choice := resp.Choices[0]
	if choice.Message.Content == nil {
		return fmt.Errorf("no content in choice")
	}

	fmt.Println()
	for _, msg := range msgs {
		fmt.Printf("%s: %s\n", msg.Role, *msg.Content)
	}
	fmt.Printf("Assistant: %s\n", *choice.Message.Content)
}
```

## Advanced Configuration

The client can be customized with various options:

```go
// With a custom zap logger
zapLogger, _ := zap.NewProduction()
client, err := openai.New(cfg, openai.WithZapLogger(zapLogger))

// With a custom URL (for proxies or compatible APIs)
client, err := openai.New(cfg, openai.WithUrl("https://your-proxy.example.com"))

// With a custom HTTP client
customClient := retryablehttp.NewClient()
customClient.RetryMax = 5
client, err := openai.New(cfg, openai.WithRetryableHttpClient(customClient))
```

## Function/Tool Calling

The package provides built-in support for OpenAI's function calling feature:

```go
// Define tools
weatherTool := openai.FunctionTool{
    Function: openai.Function{
        Name:        "get_weather",
        Description: "Get the current weather in a location",
        Parameters: json.RawMessage(`{
            "type": "object",
            "properties": {
                "location": {
                    "type": "string",
                    "description": "The city and state"
                },
                "unit": {
                    "type": "string",
                    "enum": ["celsius", "fahrenheit"]
                }
            },
            "required": ["location"]
        }`),
    },
}

// Create chat completion with tool support
opts := func(options *openai.ChatCompletionOptions) {
    options.Tools = []openai.Tool{weatherTool}
    options.ToolChoice = "auto"
}

resp, err := client.ChatCompletionCreate(context.Background(), messages, opts)
```

## Integration with GoLLuM

This package works seamlessly with other GoLLuM modules:

- Use with [Scrolls](../scrolls) for prompt templating and management
- Use with [Ringchain](../ringchain) for building complex agent workflows