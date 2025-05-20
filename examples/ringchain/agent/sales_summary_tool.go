package agent

import (
	"context"
	"maps"

	"github.com/dskart/gollum/openai"
	"github.com/dskart/gollum/ringchain"
	"go.uber.org/zap"
)

type SalesSummaryTool struct {
	graph        *ringchain.Graph
	lastNodeHash string
}

func NewSalesSummaryTool(llm openai.OpenAi) (*SalesSummaryTool, error) {
	g := ringchain.NewGraph()

	salesDataRetrievalNode, err := NewSalesDataRetrievalNode()
	if err != nil {
		return nil, err
	}

	salesDataSumNode, err := NewSalesDataSumNode()
	if err != nil {
		return nil, err
	}

	salesTotalSummarizerNode, err := NewSalesSummarizerNode(llm)
	if err != nil {
		return nil, err
	}

	if err := g.AddNode(salesDataRetrievalNode.Name(), salesDataRetrievalNode); err != nil {
		return nil, err
	}
	if err := g.AddNode(salesDataSumNode.Name(), salesDataSumNode); err != nil {
		return nil, err
	}
	if err := g.AddNode(salesTotalSummarizerNode.Name(), salesTotalSummarizerNode); err != nil {
		return nil, err
	}

	if err := g.AddEdge(salesDataRetrievalNode.Name(), salesDataSumNode.Name()); err != nil {
		return nil, err
	}
	if err := g.AddEdge(salesDataSumNode.Name(), salesTotalSummarizerNode.Name()); err != nil {
		return nil, err
	}

	return &SalesSummaryTool{
		graph:        g,
		lastNodeHash: salesTotalSummarizerNode.Name(),
	}, nil
}

const salesSummaryToolOpenAiFnName = "report_modification_tool"

var (
	description          = "Gets a sales summary for a type of product."
	salesSummaryOpenAiFn = openai.Function{
		Name:        salesSummaryToolOpenAiFnName,
		Description: &description,
		Parameters: openai.Parameters{
			Type: openai.ObjectParameterType,
			Properties: map[string]openai.Property{
				"product_type": {
					Type:        openai.StringPropertyType,
					Enum:        []string{"cars", "planes"},
					Description: "The type of product to get the sales summary for.",
				},
			},
			Required: []string{"product_type"},
		},
	}
)

func (n *SalesSummaryTool) OpenAiTool() openai.Tool {
	return openai.Tool{
		Type:     openai.FunctionToolType,
		Function: salesSummaryOpenAiFn,
	}
}

func (n *SalesSummaryTool) FunctionName() string {
	return salesSummaryToolOpenAiFnName
}

func (n *SalesSummaryTool) ToolName() string {
	return "SalesSummaryTool"
}

func (n *SalesSummaryTool) Name() string {
	return "SalesSummaryTool"
}

func (n *SalesSummaryTool) Description() string {
	return description
}

func (n *SalesSummaryTool) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	results := make(map[string]any)
	res, err := n.graph.Execute(ctx, logger, args)
	if err != nil {
		return nil, err
	}
	maps.Copy(results, res[n.lastNodeHash])
	return results, nil
}
