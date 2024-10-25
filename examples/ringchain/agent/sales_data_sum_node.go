package agent

import (
	"context"
	"maps"

	"go.uber.org/zap"
)

type SalesDataSumNode struct {
}

func NewSalesDataSumNode() (*SalesDataSumNode, error) {
	return &SalesDataSumNode{}, nil
}

func (n *SalesDataSumNode) Name() string {
	return "SalesDataSumNode"
}

func (n *SalesDataSumNode) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	results := make(map[string]any)
	entries := args["entries"].([]Entry)

	personTotalMap := make(map[string]int64)
	for _, entry := range entries {
		personTotalMap[entry.Name] += entry.SalePrice
	}

	maps.Copy(results, args)
	results["sales_sum"] = personTotalMap
	return results, nil
}
