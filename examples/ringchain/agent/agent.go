package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dskart/gollum/openai"
	"github.com/dskart/gollum/ringchain"
	"github.com/dskart/gollum/scrolls"
	"go.uber.org/zap"
)

type MyAgent struct {
	scroll *scrolls.Scroll
	llm    openai.OpenAi
	tools  []ringchain.Tool
}

func NewMyAgent(llm openai.OpenAi) (*MyAgent, error) {
	salesSummaryTool, err := NewSalesSummaryTool(llm)
	if err != nil {
		return nil, err
	}

	tools := []ringchain.Tool{
		salesSummaryTool,
	}

	scrollText := `
[[#system~]]
Today's date is the {{.current_date}}.
You are an expert at Sales Data Analysis, and you are working with {{.domains}} data.
Your task is to help the user analyse their product sales data.

To achieve your goal, you have access to the following tools through function calls:

{{ range $i, $value := .tools }}
{{$value.ToolName}}: {{$value.Description}}
{{ end }}
[[~/system]]

[[#user~]]
{{.question}}
[[~/user]]
`
	scroll := scrolls.New(scrollText, llm)

	return &MyAgent{
		scroll: scroll,
		llm:    llm,
		tools:  tools,
	}, nil

}

func (c *MyAgent) Name() string {

	return "MyAgentNode"
}

func (c *MyAgent) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	results := make(map[string]any)

	question := args["question"].(string)

	msgs, err := c.scroll.ParseBlocks(map[string]any{
		"tools":        c.tools,
		"current_date": time.Now().Format("2006-01-02"),
		"domain":       []string{"cars", "planes"},
		"question":     question,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.llm.ChatCompletionCreate(
		ctx,
		msgs,
		openai.WithTools(ringchain.OpenAiFunctions(c.tools)),
		openai.WithN(1),
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("zero choices returned")
	}

	choice := resp.Choices[0]
	if choice.IsAssistantMessage() {
		content := choice.Message.Content
		results["assistant_msg"] = *content
	} else if choice.IsToolCall() {
		selectedFnName := choice.Message.ToolCalls[0].Function.Name
		selectedFnArgs := choice.Message.ToolCalls[0].Function.Arguments

		tool, ok := ringchain.SelectTool(c.tools, selectedFnName)
		if !ok {
			return nil, fmt.Errorf("cannot find tool: %s", selectedFnName)
		}

		var toolArgs map[string]any
		err := json.Unmarshal([]byte(selectedFnArgs), &toolArgs)
		if err != nil {
			return nil, err
		}

		toolResults, err := tool.Run(ctx, logger, toolArgs)
		if err != nil {
			return nil, err
		}

		summary := toolResults["summary"]
		results["assistant_msg"] = summary
	}

	return results, nil
}
