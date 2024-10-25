package openai

// https://platform.openai.com/docs/api-reference/chat/object
type ChatCompletionObject struct {
	Id                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Object            string   `json:"object"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	FinishReason FinishReasonType      `json:"finish_reason"`
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
}

type FinishReasonType string

const (
	StopFinishReasonType          FinishReasonType = "stop"
	ToolCallsFinishReasonType     FinishReasonType = "tool_calls"
	ContentFilterFinishReasonType FinishReasonType = "content_filter"
	LengthFinishReasonType        FinishReasonType = "length"
)

func (f FinishReasonType) String() string {
	return string(f)
}

func (c Choice) IsAssistantMessage() bool {
	return c.FinishReason == StopFinishReasonType && c.Message.Content != nil
}

func (c Choice) IsToolCall() bool {
	return c.FinishReason == ToolCallsFinishReasonType && c.Message.Content == nil
}

func (c Choice) IsStop() bool {
	return c.FinishReason == StopFinishReasonType
}

func (c Choice) IsLength() bool {
	return c.FinishReason == LengthFinishReasonType
}

func (c Choice) IsContentFilter() bool {
	return c.FinishReason == ContentFilterFinishReasonType
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
