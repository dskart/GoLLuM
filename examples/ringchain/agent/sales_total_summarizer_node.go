package agent

import (
	"context"
	"fmt"
	"maps"

	"github.com/dskart/gollum/openai"
	"github.com/dskart/gollum/scrolls"
	"go.uber.org/zap"
)

type SalesSummarizerNode struct {
	scroll *scrolls.Scroll
}

func NewSalesSummarizerNode(llm openai.OpenAi) (*SalesSummarizerNode, error) {
	scrollText := `
[[#system~]]
You are a data summarization expert.
[[~/system]]

[[#user~]]
Here are the sales total for each person:

Result:
{{.sales_sum}}

Can you summarize the data for the given product sales {{.product_type}}
[[~/user]]

[[#assistant~]]
{"action": "gen", "output_name": "summary", "temperature": 0, "max_tokens": 300}
[[~/assistant]]
`

	scroll := scrolls.New(scrollText, llm)
	return &SalesSummarizerNode{
		scroll: scroll,
	}, nil
}

func (n *SalesSummarizerNode) Name() string {
	return "SalesSummarizerNode"
}

func (n *SalesSummarizerNode) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	results := make(map[string]any)

	salesSum := args["sales_sum"].(map[string]int64)
	productType := args["product_type"].(string)

	_, res, err := n.scroll.Execute(ctx, map[string]any{
		"product_type": productType,
		"sales_sum":    salesSum,
	},
	)
	if err != nil {
		return nil, err
	}
	summary, ok := res["summary"]
	if !ok {
		return nil, fmt.Errorf("could not find summary in scroll result")
	}

	maps.Copy(results, args)
	results["summary"] = summary
	return results, nil
}
