package ringchain

import (
	"context"

	"go.uber.org/zap"
)

type Node interface {
	Name() string

	// Run executes the node with the given arguments and returns the result.
	// If you want to copy the args map into the result, use maps.Clone(args).
	Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error)
}

type Edge struct {
	SourceHash string
	TargetHash string
}
