package ringchain

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"go.uber.org/zap"
)

type Graph struct {
	store Store
}

func NewGraph() *Graph {
	return &Graph{
		store: newMemoryStore(),
	}
}

func (g *Graph) AddNode(name string, node Node) error {
	hash := name
	return g.store.AddNode(hash, node)
}

func (g *Graph) Node(hash string) (Node, error) {
	node, err := g.store.Node(hash)
	return node, err
}

func (g *Graph) AddEdge(sourceHash, targetHash string) error {
	_, err := g.store.Node(sourceHash)
	if err != nil {
		return fmt.Errorf("source node %v: %w", sourceHash, err)
	}

	_, err = g.store.Node(targetHash)
	if err != nil {
		return fmt.Errorf("target node %v: %w", targetHash, err)
	}

	if _, err := g.Edge(sourceHash, targetHash); !errors.Is(err, ErrEdgeNotFound) {
		return ErrEdgeAlreadyExists
	}

	createsCycle, err := g.createsCycle(sourceHash, targetHash)
	if err != nil {
		return fmt.Errorf("check for cycles: %w", err)
	} else if createsCycle {
		return ErrEdgeCreatesCycle
	}

	edge := Edge{
		SourceHash: sourceHash,
		TargetHash: targetHash,
	}

	return g.store.AddEdge(sourceHash, targetHash, edge)
}

func (g *Graph) Edge(sourceHash string, targetHash string) (Edge, error) {
	edge, err := g.store.Edge(sourceHash, targetHash)
	if err != nil {
		return Edge{}, err
	}

	return edge, nil
}

func (g *Graph) SuccessorMap() (map[string]map[string]Edge, error) {
	vertices, err := g.store.ListNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	edges, err := g.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	m := make(map[string]map[string]Edge, len(vertices))

	for _, node := range vertices {
		m[node] = make(map[string]Edge)
	}

	for _, edge := range edges {
		m[edge.SourceHash][edge.TargetHash] = edge
	}

	return m, nil
}

func (g *Graph) PredecessorMap() (map[string]map[string]Edge, error) {
	vertices, err := g.store.ListNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to list vertices: %w", err)
	}

	edges, err := g.store.ListEdges()
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	m := make(map[string]map[string]Edge, len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[string]Edge)
	}

	for _, edge := range edges {
		if _, ok := m[edge.TargetHash]; !ok {
			m[edge.TargetHash] = make(map[string]Edge)
		}
		m[edge.TargetHash][edge.SourceHash] = edge
	}

	return m, nil
}

func (g *Graph) createsCycle(source, target string) (bool, error) {
	return g.store.CreatesCycle(source, target)
}

type GraphExecuteOptions struct {
	NumWorkers int
}

func WithNumWorkers(n int) func(*GraphExecuteOptions) {
	return func(opts *GraphExecuteOptions) {
		opts.NumWorkers = n
	}
}

// Execute executes the graph starting from the given node.
// You need to call Init before calling Execute.
// This method is safe to call concurrently but might break if the graph is modified while executing.
func (g *Graph) Execute(ctx context.Context, logger *zap.Logger, args map[string]any, opts ...func(*GraphExecuteOptions)) (map[string]map[string]any, error) {
	options := GraphExecuteOptions{
		NumWorkers: 10,
	}
	for _, opt := range opts {
		opt(&options)
	}

	predecessorMap, err := g.PredecessorMap()
	if err != nil {
		return nil, err
	}
	nodeCount, err := g.store.NodeCount()
	if err != nil {
		return nil, err
	}

	workerPool := newWorkerPool(options.NumWorkers, nodeCount)
	workerPool.Run(ctx)

	// add all nodes without predecessors to the worker pool
	for nodeHash, predecessors := range predecessorMap {
		if len(predecessors) == 0 {
			node, err := g.Node(nodeHash)
			if err != nil {
				return nil, err
			}
			workerPool.AddNodeJob(logger, nodeHash, node, args)
			delete(predecessorMap, nodeHash)
		}
	}

	var nodeResults = make(map[string]map[string]any, nodeCount)
	var nodeInputMap = make(map[string]map[string]any, nodeCount)

	for workerPool.NumRunningNodes() > 0 {
		runningNode := workerPool.PopRunningNode()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-workerPool.errChan:
			return nil, err
		case res := <-runningNode.Result:

			nodeResults[runningNode.NodeHash] = res
			for nodeHash, predecessors := range predecessorMap {
				// check for predecessors of the finished node
				if _, ok := predecessors[runningNode.NodeHash]; ok {
					if _, ok := nodeInputMap[nodeHash]; !ok {
						nodeInputMap[nodeHash] = make(map[string]any)
					}
					// Save running node result in the next node input map
					maps.Copy(nodeInputMap[nodeHash], res)
					delete(predecessors, runningNode.NodeHash)
				}

				if len(predecessors) == 0 {
					node, err := g.Node(nodeHash)
					if err != nil {
						return nil, err
					}
					input := nodeInputMap[nodeHash]

					workerPool.AddNodeJob(logger, nodeHash, node, input)
					delete(predecessorMap, nodeHash)
				}
			}
		default:
			// no result yet, put it back in the queue
			workerPool.AppendRunningNode(runningNode)
		}
	}

	workerPool.Stop()
	if err := workerPool.Wait(); err != nil {
		return nil, err
	}

	return nodeResults, nil
}
