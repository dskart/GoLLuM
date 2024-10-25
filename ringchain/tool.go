package ringchain

import (
	"context"

	"github.com/dskart/gollum/openai"
	"go.uber.org/zap"
)

type Tool interface {
	OpenAiTool() openai.Tool
	FunctionName() string
	Description() string
	ToolName() string
	Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error)
}

func OpenAiFunctions(tools []Tool) []openai.Tool {
	ret := make([]openai.Tool, len(tools))
	for i, tool := range tools {
		ret[i] = tool.OpenAiTool()
	}
	return ret
}

func SelectTool(tools []Tool, functionName string) (Tool, bool) {
	for _, tool := range tools {
		if tool.FunctionName() == functionName {
			return tool, true
		}
	}
	return nil, false
}
