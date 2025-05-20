package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dskart/gollum/openai"
)

func run(ctx context.Context, openAiKey string, gptModel string) error {
	openaiCfg := openai.Config{
		OpenAiKey: openAiKey,
		GptModel:  gptModel,
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
