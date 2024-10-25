package scrolls

import (
	"context"
	"fmt"

	"github.com/dskart/gollum/openai"
)

type AssistantBody struct {
	Action           string              `json:"action"`
	OutputName       string              `json:"output_name"`
	FrequencyPenalty *float64            `json:"frequency_penalty,omitempty"`
	LogitBias        *map[string]float64 `json:"logit_bias,omitempty"`
	MaxToken         *int                `json:"max_tokens,omitempty"`
	N                *int                `json:"n,omitempty"`
	PresencePenalty  *float64            `json:"presence_penalty,omitempty"`
	Stop             *[]string           `json:"stop,omitempty"`
	Temperature      *float64            `json:"temperature,omitempty"`
	TopP             *float64            `json:"top_p,omitempty"`
	User             *string             `json:"user,omitempty"`
}

type ActionType string

const (
	GenActionType string = "gen"
)

func (a AssistantBody) OpenAiPromptOptions() []func(*openai.ChatCompletionOptions) {
	ret := make([]func(*openai.ChatCompletionOptions), 0, 9)
	if a.FrequencyPenalty != nil {
		ret = append(ret, openai.WithFrequencyPenalty(*a.FrequencyPenalty))
	}

	if a.LogitBias != nil {
		ret = append(ret, openai.WithLogitBias(*a.LogitBias))
	}

	if a.MaxToken != nil {
		ret = append(ret, openai.WithMaxToken(*a.MaxToken))
	}

	if a.N != nil {
		ret = append(ret, openai.WithN(*a.N))
	}

	if a.PresencePenalty != nil {
		ret = append(ret, openai.WithPresencyPenalty(*a.PresencePenalty))
	}

	if a.Stop != nil {
		ret = append(ret, openai.WithStop(*a.Stop))
	}

	if a.Temperature != nil {
		ret = append(ret, openai.WithTemperature(*a.Temperature))
	}

	if a.TopP != nil {
		ret = append(ret, openai.WithTopP(*a.TopP))
	}

	if a.User != nil {
		ret = append(ret, openai.WithUser(*a.User))
	}

	return ret
}

func promptOpenAi(ctx context.Context, llm openai.OpenAi, msgs []openai.Message, opts ...func(*openai.ChatCompletionOptions)) (string, error) {
	resp, err := llm.ChatCompletionCreate(ctx, msgs, opts...)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	choice := resp.Choices[0]
	if choice.Message.Content == nil {
		return "", fmt.Errorf("no content in choice")
	}
	choiceContent := *choice.Message.Content
	return choiceContent, nil
}
