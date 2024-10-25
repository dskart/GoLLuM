package scrolls

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
	"time"

	"github.com/dskart/gollum/openai"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TEST_MODEL = "gpt-4"
	TEST_KEY   = "foo"
)

func testHttpClient() *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 0
	retryClient.RetryWaitMin = 0 * time.Second
	retryClient.RetryWaitMax = 0 * time.Second
	retryClient.Logger = nil
	return retryClient
}

func newTestServer(t *testing.T, expectedReqBody openai.ChatCompletionRequestBody, resp openai.ChatCompletionObject) *httptest.Server {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		reqData, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var reqBody openai.ChatCompletionRequestBody
		err = json.Unmarshal(reqData, &reqBody)
		require.NoError(t, err)

		assert.Equal(t, expectedReqBody, reqBody)

		rawBody, err := json.Marshal(&resp)
		require.NoError(t, err)
		w.Write(rawBody)
	}))
	return svr
}

var testTemplate string = `
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

func toPointer[T any](t T) *T {
	return &t
}

func TestScrolls(t *testing.T) {
	ctx := context.Background()
	content := "I'm Sorry Dave, I'm Afraid I Can't Do That"
	resp := openai.ChatCompletionObject{
		Choices: []openai.Choice{
			openai.Choice{
				Message: openai.ChatCompletionMessage{
					Role:    openai.AssistantRoleType,
					Content: &content,
				},
			},
		},
	}
	msgs := []openai.Message{
		{
			Role:    openai.SystemRoleType,
			Content: toPointer("You are a helpful and terse assistant."),
		},
		{
			Role:    openai.UserRoleType,
			Content: toPointer("I want a response to the following question: What is the meaning of life?"),
		},
	}
	expectedBody := openai.ChatCompletionRequestBody{
		Model:       TEST_MODEL,
		Messages:    msgs,
		Temperature: toPointer(0.0),
		MaxToken:    toPointer(300),
	}

	svr := newTestServer(t, expectedBody, resp)
	defer svr.Close()

	openAi, err := openai.New(openai.Config{
		OpenAiKey: TEST_KEY,
		GptModel:  TEST_MODEL,
	}, openai.WithRetryableHttpClient(testHttpClient()), openai.WithUrl(svr.URL))
	require.NoError(t, err)

	scroll := New(testTemplate, openAi)

	results, _, err := scroll.Execute(ctx, map[string]any{
		"query": "What is the meaning of life?",
	})
	require.NoError(t, err)

	assert.Equal(t, []openai.Message{
		{
			Role:    openai.SystemRoleType,
			Content: toPointer("You are a helpful and terse assistant."),
		},
		{
			Role:    openai.UserRoleType,
			Content: toPointer("I want a response to the following question: What is the meaning of life?"),
		},
		{
			Role:    openai.AssistantRoleType,
			Content: &content,
		},
	}, results)
}

var testTemplate2 string = `
[[#system~]]
You are a helpful and terse assistant.
[[~/system]]


{{ range $i, $value := .questions }}
[[#user~]]
I want a response to the following question: {{question_str $value}}
[[~/user]]
{{ end }}

{{if .add_tip}}
[[#user~]]
TIP:
You can ask me for a tip by saying "tip".
[[~/user]]

[[#system~]]
Thanks for the tip
[[~/system]]
{{end}}

[[#assistant~]]
{"action": "gen", "output_name": "response", "temperature": 0, "max_tokens": 300}
[[~/assistant]]
`

func TestParseBlocks(t *testing.T) {
	openAi, err := openai.New(openai.Config{
		OpenAiKey: TEST_KEY,
		GptModel:  TEST_MODEL,
	}, openai.WithRetryableHttpClient(testHttpClient()), openai.WithUrl(""))
	require.NoError(t, err)

	scroll := New(testTemplate2, openAi, WithFuncMap(template.FuncMap{
		"question_str": func(question string) string {
			return question + "?"
		},
	}))

	questions := []string{
		"foo",
		"bar",
	}

	blocks, err := scroll.ParseBlocks(map[string]any{
		"questions": questions,
	})
	require.NoError(t, err)
	expected := []openai.Message{
		{Role: openai.SystemRoleType, Content: toPointer("You are a helpful and terse assistant.")},
		{Role: openai.UserRoleType, Content: toPointer("I want a response to the following question: foo?")},
		{Role: openai.UserRoleType, Content: toPointer("I want a response to the following question: bar?")},
		{Role: openai.AssistantRoleType, Content: toPointer(`{"action": "gen", "output_name": "response", "temperature": 0, "max_tokens": 300}`)},
	}
	assert.Equal(t, expected, blocks)
}
