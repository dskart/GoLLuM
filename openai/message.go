package openai

type RoleType string

const (
	SystemRoleType    RoleType = "system"
	UserRoleType      RoleType = "user"
	AssistantRoleType RoleType = "assistant"
	ToolRoleType      RoleType = "tool"
)

type Message struct {
	Content    *string    `json:"content"`
	Role       RoleType   `json:"role"`
	Name       *string    `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallId string     `json:"tool_call_id"`
}

type ChatCompletionMessage struct {
	Content   *string    `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
	Role      RoleType   `json:"role"`
}
