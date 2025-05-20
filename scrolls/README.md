# 📜 Scrolls

A powerful templating system for managing LLM prompts and conversations in Go applications. Scrolls makes it easy to create, maintain, and execute complex conversation templates with OpenAI models.

## Features

- **Template-Based Prompts**: Define complex prompt templates with Go's template syntax
- **Role-Based Messages**: Easily create system, user, and assistant messages
- **Chain of Thought**: Support for multi-turn conversations
- **Dynamic Content**: Inject variables into your templates at runtime
- **Stateful Execution**: Track and manage conversation state

## Installation

```bash
go get github.com/dskart/gollum/scrolls
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dskart/gollum/openai"
	"github.com/dskart/gollum/scrolls"
)

func main() {
	// Set up the OpenAI client
	cfg := openai.Config{
		GptModel:  "gpt-4",
		OpenAiKey: os.Getenv("OPENAI_API_KEY"),
	}
	client, _ := openai.New(cfg)

	// Define a template with role-based blocks
	template := `
[[#system~]]
You are a helpful assistant that provides information about {{.topic}}.
[[~/system]]

[[#user~]]
Tell me 3 interesting facts about {{.topic}}.
[[~/user]]
`

	// Create a new scroll with the template
	scroll := scrolls.New(template, client)

	// Execute the template with arguments
	messages, outputs, err := scroll.Execute(context.Background(), map[string]any{
		"topic": "quantum computing",
	})
	
	if err != nil {
		fmt.Printf("Error executing scroll: %v\n", err)
		return
	}

	// Print the conversation
	for _, msg := range messages {
		fmt.Printf("%s: %s\n", msg.Role, *msg.Content)
	}
}
```

## Assistant Actions

Scrolls supports special "assistant" blocks that can generate content dynamically:

```go
template := `
[[#system~]]
You are a helpful assistant.
[[~/system]]

[[#user~]]
Tell me about {{.topic}}.
[[~/user]]

[[#assistant~]]
{
  "type": "openai",
  "output_name": "response"
}
[[~/assistant]]
`

scroll := scrolls.New(template, client)
_, outputs, _ := scroll.Execute(context.Background(), map[string]any{
	"topic": "artificial intelligence",
})

// Access the generated content
aiResponse := outputs["response"]
fmt.Println(aiResponse)
```

## Custom Template Functions

You can extend Scrolls with custom template functions:

```go
funcMap := template.FuncMap{
	"uppercase": strings.ToUpper,
	"currentDate": func() string {
		return time.Now().Format("2006-01-02")
	},
}

scroll := scrolls.New(template, client, scrolls.WithFuncMap(funcMap))
```

## Integration with GoLLuM

This package works seamlessly with other GoLLuM modules:

- Use with [OpenAI](../openai) for API communication
- Use with [Ringchain](../ringchain) for building complex agent workflows