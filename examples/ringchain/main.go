package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dskart/gollum/examples/ringchain/agent"
	"github.com/dskart/gollum/openai"
	"go.uber.org/zap"
)

func run(ctx context.Context, logger *zap.Logger, openAiKey string, gptModel string) error {
	llm, err := openai.New(openai.Config{
		OpenAiKey: openAiKey,
		GptModel:  gptModel,
	}, openai.WithZapLogger(logger))
	if err != nil {
		return fmt.Errorf("failed to initialize openai: %w", err)
	}

	agent, err := agent.NewMyAgent(llm)
	if err != nil {
		return fmt.Errorf("failed to initialize agent: %w", err)
	}

	args := map[string]any{
		"question": "Hello, how are you doing",
	}

	results, err := agent.Run(ctx, logger, args)
	if err != nil {
		return fmt.Errorf("failed to run agent: %w", err)
	}

	for k, v := range results {
		fmt.Printf("%s: %v\n", k, v)
	}

	return nil
}

func main() {
	ctx := context.Background()
	openAiKey := os.Getenv("OPENAI_API_KEY")
	gptModel := os.Getenv("GPT_MODEL")

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if err := run(ctx, logger, openAiKey, gptModel); err != nil {
		panic(err)
	}
}
