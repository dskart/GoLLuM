package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func newTestServer(t *testing.T, expectedReqBody ChatCompletionRequestBody) *httptest.Server {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			if err != nil {
				t.Errorf("Failed to close request body: %v", err)
			}
		}()
		reqData, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var reqBody ChatCompletionRequestBody
		err = json.Unmarshal(reqData, &reqBody)
		require.NoError(t, err)

		assert.Equal(t, expectedReqBody, reqBody)

		resp := ChatCompletionObject{}
		rawBody, err := json.Marshal(&resp)
		require.NoError(t, err)
		_, err = w.Write(rawBody)
		require.NoError(t, err)
	}))
	return svr
}

func strPointer(s string) *string {
	return &s
}

func ChatCompletionCreate(t *testing.T) {
	ctx := context.Background()

	msgs := []Message{
		{
			Role:    UserRoleType,
			Content: strPointer("Open The pod bay doors, HAL."),
		},
		{
			Role:    AssistantRoleType,
			Content: strPointer("I'm Sorry Dave, I'm Afraid I Can't Do That"),
		},
	}
	expected := ChatCompletionRequestBody{
		Model:    TEST_MODEL,
		Messages: msgs,
	}

	svr := newTestServer(t, expected)
	defer svr.Close()

	openAi, err := New(Config{
		OpenAiKey: TEST_KEY,
		GptModel:  TEST_MODEL,
	}, WithRetryableHttpClient(testHttpClient()), WithUrl(svr.URL))
	require.NoError(t, err)

	_, err = openAi.ChatCompletionCreate(ctx, msgs)
	require.NoError(t, err)
}

func TestChatCompletionCreateWithOptions(t *testing.T) {
	ctx := context.Background()
	msgs := []Message{
		{Role: UserRoleType, Content: strPointer("What is the meaning of life?")},
	}

	expectedReqBody := func(opt func(*ChatCompletionOptions)) ChatCompletionRequestBody {
		options := ChatCompletionOptions{}
		opt(&options)

		return ChatCompletionRequestBody{
			Messages:         msgs,
			Model:            TEST_MODEL,
			FrequencyPenalty: options.frequencyPenalty,
			Tools:            options.tools,
			ToolChoice:       options.toolChoice,
			LogitBias:        options.logitBias,
			MaxToken:         options.maxToken,
			N:                options.n,
			PresencyPenalty:  options.presencyPenalty,
			Stop:             options.stop,
			Temperature:      options.temperature,
			TopP:             options.topP,
			User:             options.user,
		}
	}

	testCases := []struct {
		fieldName string
		option    func(*ChatCompletionOptions)
	}{
		{fieldName: "FrequencyPenalty", option: WithFrequencyPenalty(0.5)},
		{fieldName: "Tools", option: WithTools([]Tool{{Type: FunctionToolType, Function: Function{Name: "foo"}}})},
		{fieldName: "ToolChoice", option: WithToolChoice("foo")},
		{fieldName: "LogitBias", option: WithLogitBias(map[string]float64{"foo": 0.5})},
		{fieldName: "MaxToken", option: WithMaxToken(10)},
		{fieldName: "N", option: WithN(10)},
		{fieldName: "PresencyPenalty", option: WithPresencyPenalty(0.5)},
		{fieldName: "Stop", option: WithStop([]string{"foo"})},
		{fieldName: "Temperature", option: WithTemperature(0.5)},
		{fieldName: "TopP", option: WithTopP(0.5)},
		{fieldName: "User", option: WithUser("foo")},
	}
	for _, tc := range testCases {
		t.Run(tc.fieldName, func(t *testing.T) {
			msgs := []Message{
				{Role: UserRoleType, Content: strPointer("What is the meaning of life?")},
			}

			expected := expectedReqBody(tc.option)
			svr := newTestServer(t, expected)
			defer svr.Close()

			openAi, err := New(Config{
				OpenAiKey: TEST_KEY,
				GptModel:  TEST_MODEL,
			}, WithRetryableHttpClient(testHttpClient()), WithUrl(svr.URL))
			require.NoError(t, err)

			resp, err := openAi.ChatCompletionCreate(ctx, msgs, tc.option)
			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}
