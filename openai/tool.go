package openai

type ToolCall struct {
	Id       string       `json:"id"`
	Type     ToolType     `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Tool struct {
	Type     ToolType `json:"type"`
	Function Function `json:"function"`
}

type ToolType string

const (
	FunctionToolType ToolType = "function"
)

type Function struct {
	Description *string    `json:"description,omitempty"`
	Name        string     `json:"name"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       ParameterType       `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

type ParameterType string

const (
	ObjectParameterType ParameterType = "object"
)

type Property struct {
	Type        PropertyType `json:"type"`
	Description string       `json:"description"`
	// Enum is optional
	Enum []string `json:"enum,omitempty"`
	// Items is only valid for array types
	Items *Items `json:"items,omitempty"`
}

type PropertyType string

const (
	StringPropertyType  PropertyType = "string"
	IntegerPropertyType PropertyType = "integer"
	ArrayPropertyType   PropertyType = "array"
)

type Items struct {
	Type PropertyType `json:"type"`
}
