package ringchain

import "errors"

var (
	ErrNodeNotFound      = errors.New("vertex not found")
	ErrNodeAlreadyExists = errors.New("vertex already exists")
	ErrEdgeNotFound      = errors.New("edge not found")
	ErrEdgeAlreadyExists = errors.New("edge already exists")
	ErrEdgeCreatesCycle  = errors.New("edge would create a cycle")
	ErrNodeHasEdges      = errors.New("vertex has edges")
	ErrGraphNotInit      = errors.New("graph not Init()")
)
