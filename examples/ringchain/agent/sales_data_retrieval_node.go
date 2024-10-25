package agent

import (
	"context"
	"fmt"
	"maps"

	"go.uber.org/zap"
)

type SalesDataRetrievalNode struct {
}

type Entry struct {
	Name      string
	SalePrice int64
}

var db map[string][]Entry = map[string][]Entry{
	"cars": {
		Entry{Name: "Jeffrey", SalePrice: 10},
		Entry{Name: "Jeffrey", SalePrice: 20},
		Entry{Name: "Alice", SalePrice: 30},
		Entry{Name: "Alice", SalePrice: 40},
	},
	"planes": {
		Entry{Name: "Jeffrey", SalePrice: 100},
		Entry{Name: "Jeffrey", SalePrice: 200},
		Entry{Name: "Alice", SalePrice: 300},
		Entry{Name: "Alice", SalePrice: 400},
	},
}

func NewSalesDataRetrievalNode() (*SalesDataRetrievalNode, error) {
	return &SalesDataRetrievalNode{}, nil
}

func (n *SalesDataRetrievalNode) Name() string {
	return "SalesDataRetrievalNode"
}

func (n *SalesDataRetrievalNode) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	results := make(map[string]any)
	productType := args["product_type"].(string)

	entries, ok := db[productType]
	if !ok {
		return nil, fmt.Errorf("product type %s not found", productType)
	}

	maps.Copy(results, args)
	results["entries"] = entries
	return results, nil
}
