package openai

import "context"

type ChatCompletionOptions struct {
	frequencyPenalty *float64
	logitBias        *map[string]float64
	maxToken         *int
	n                *int
	presencyPenalty  *float64
	responseFormat   *ResponseFormat
	stop             *[]string
	stream           *bool
	temperature      *float64
	topP             *float64
	tools            *[]Tool
	toolChoice       *string
	user             *string
}

func WithFrequencyPenalty(frequencyPenalty float64) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.frequencyPenalty = &frequencyPenalty
	}
}

func WithLogitBias(logitBias map[string]float64) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.logitBias = &logitBias
	}
}

func WithMaxToken(maxToken int) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.maxToken = &maxToken
	}
}

func WithN(n int) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.n = &n
	}
}

func WithPresencyPenalty(presencyPenalty float64) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.presencyPenalty = &presencyPenalty
	}
}

func WithResponseFormat(responseFormat ResponseFormat) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.responseFormat = &responseFormat
	}
}

func WithStop(stop []string) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.stop = &stop
	}
}

func WithStream(stream bool) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.stream = &stream
	}
}

func WithTemperature(temperature float64) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.temperature = &temperature
	}
}

func WithTopP(topP float64) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.topP = &topP
	}
}

func WithTools(tools []Tool) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.tools = &tools
	}
}

func WithToolChoice(toolChoice string) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.toolChoice = &toolChoice
	}
}

func WithUser(user string) func(*ChatCompletionOptions) {
	return func(opts *ChatCompletionOptions) {
		opts.user = &user
	}
}

func (o *OpenAiImpl) ChatCompletionCreate(ctx context.Context, messages []Message, opts ...func(*ChatCompletionOptions)) (ChatCompletionObject, error) {
	options := ChatCompletionOptions{}
	for _, o := range opts {
		o(&options)
	}

	reqBody := ChatCompletionRequestBody{
		Messages:         messages,
		Model:            o.gptModel,
		FrequencyPenalty: options.frequencyPenalty,
		LogitBias:        options.logitBias,
		MaxToken:         options.maxToken,
		N:                options.n,
		PresencyPenalty:  options.presencyPenalty,
		ResponseFormat:   options.responseFormat,
		Stop:             options.stop,
		Stream:           options.stream,
		Temperature:      options.temperature,
		TopP:             options.topP,
		Tools:            options.tools,
		ToolChoice:       options.toolChoice,
		User:             options.user,
	}

	resp, err := request[ChatCompletionRequestBody, ChatCompletionObject](ctx, o.httpClient, o.getChatCompletionUrl(), o.apiKey, reqBody)
	if err != nil {
		return ChatCompletionObject{}, err
	}

	return *resp, nil
}

// https://platform.openai.com/docs/api-reference/chat/create
type ChatCompletionRequestBody struct {
	Messages         []Message           `json:"messages"`
	Model            string              `json:"model"`
	FrequencyPenalty *float64            `json:"frequency_penalty,omitempty"`
	LogitBias        *map[string]float64 `json:"logit_bias,omitempty"`
	MaxToken         *int                `json:"max_tokens,omitempty"`
	N                *int                `json:"n,omitempty"`
	PresencyPenalty  *float64            `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat     `json:"response_format,omitempty"`
	Seed             *int                `json:"seed,omitempty"`
	Stop             *[]string           `json:"stop,omitempty"`
	Stream           *bool               `json:"stream,omitempty"`
	Temperature      *float64            `json:"temperature,omitempty"`
	TopP             *float64            `json:"top_p,omitempty"`
	Tools            *[]Tool             `json:"tools,omitempty"`
	ToolChoice       *string             `json:"tool_choice,omitempty"`
	User             *string             `json:"user,omitempty"`
}

type ResponseFormat struct {
	Type ResponseFormatType `json:"type"`
}

type ResponseFormatType string

const (
	TextResponseFormatType       ResponseFormatType = "text"
	JsonObjectResponseFormatType ResponseFormatType = "json_object"
)
