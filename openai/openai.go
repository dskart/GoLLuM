package openai

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

// https://platform.openai.com/docs/guides/gpt/function-calling

type OpenAi interface {
	ChatCompletionCreate(ctx context.Context, messages []Message, opts ...func(*ChatCompletionOptions)) (ChatCompletionObject, error)
}

type OpenAiImpl struct {
	cfg        Config
	gptModel   string
	apiKey     string
	httpClient *retryablehttp.Client
	openAiUrl  url.URL
}

type OpenAiOptions struct {
	logger              interface{}
	retryableHttpClient *retryablehttp.Client
	url                 *string
}

func WithLogger(logger *log.Logger) func(*OpenAiOptions) {
	return func(opts *OpenAiOptions) {
		opts.logger = logger
	}
}

func WithLeveledLogger(logger retryablehttp.LeveledLogger) func(*OpenAiOptions) {
	return func(opts *OpenAiOptions) {
		opts.logger = logger
	}
}

func WithZapLogger(logger *zap.Logger) func(*OpenAiOptions) {
	return func(opts *OpenAiOptions) {
		opts.logger = retryablehttp.LeveledLogger(&LeveledZapLogger{logger})
	}
}

func WithRetryableHttpClient(retryableHttpClient *retryablehttp.Client) func(*OpenAiOptions) {
	return func(opts *OpenAiOptions) {
		opts.retryableHttpClient = retryableHttpClient
	}
}

func WithUrl(rawUrl string) func(*OpenAiOptions) {
	return func(opts *OpenAiOptions) {
		opts.url = &rawUrl
	}
}

func New(cfg Config, opts ...func(*OpenAiOptions)) (*OpenAiImpl, error) {
	options := OpenAiOptions{}
	for _, o := range opts {
		o(&options)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	openAiUrl, err := getOpenAiUrl(options.url)
	if err != nil {
		return nil, err
	}

	var httpClient *retryablehttp.Client
	if options.retryableHttpClient != nil {
		httpClient = options.retryableHttpClient
	} else {
		httpClient = setupHttpClient(options.logger)
	}

	return &OpenAiImpl{
		cfg:        cfg,
		gptModel:   cfg.GptModel,
		httpClient: httpClient,
		apiKey:     cfg.OpenAiKey,
		openAiUrl:  *openAiUrl,
	}, nil
}

func getOpenAiUrl(rawUrl *string) (*url.URL, error) {
	var openAiUrl *url.URL
	if rawUrl != nil {
		var err error
		openAiUrl, err = url.Parse(*rawUrl)
		if err != nil {
			return nil, err
		}
	} else {
		openAiUrl = &url.URL{
			Scheme: "https",
			Host:   "api.openai.com",
		}
	}

	return openAiUrl, nil
}

func setupHttpClient(logger interface{}) *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.HTTPClient.Timeout = 30 * time.Second
	retryClient.Backoff = retryablehttp.DefaultBackoff
	retryClient.Logger = logger

	return retryClient
}

func (o *OpenAiImpl) getChatCompletionUrl() string {
	return o.openAiUrl.JoinPath("v1", "chat", "completions").String()
}
