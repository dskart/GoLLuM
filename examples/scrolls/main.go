package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dskart/gollum/openai"
	"github.com/dskart/gollum/scrolls"
)

// This is just a go template
// You can inject and call functions, do loops, do conditional blocks...
var template string = `
[[#system~]]
You are a helpful and terse assistant.
[[~/system]]

[[#user~]]
I want a response to the following question: {{.query}}
[[~/user]]

[[#assistant~]]
{"action": "gen", "output_name": "response", "temperature": 0, "max_tokens": 300}
[[~/assistant]]
`

func run(ctx context.Context, openAiKey string, gptModel string) error {
	llm, err := openai.New(openai.Config{
		OpenAiKey: openAiKey,
		GptModel:  gptModel,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize openai: %w", err)
	}

	args := map[string]any{
		"query": "What is the meaning of life?",
	}

	myScroll := scrolls.New(template, llm)

	// Parse the template and returns a slice of blocks with their content
	// Use this to debug your template
	blocks, err := myScroll.ParseBlocks(args)
	if err != nil {
		return err
	}
	for _, b := range blocks {
		fmt.Printf("%s: %s\n", b.Role, *b.Content)
	}
	fmt.Println()

	// Returns a map[string]string of the #assistant blocks outputs
	_, results, err := myScroll.Execute(ctx, args)
	if err != nil {
		return err
	}

	for k, v := range results {
		fmt.Printf("%s: %s\n", k, v)
	}

	return nil
}

func main() {
	ctx := context.Background()
	openAiKey := os.Getenv("OPENAI_API_KEY")
	gptModel := os.Getenv("GPT_MODEL")

	if err := run(ctx, openAiKey, gptModel); err != nil {
		panic(err)
	}
}
