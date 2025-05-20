# Ringchain

A powerful graph-based framework for building and executing agent workflows concurrently in Go. Ringchain makes it easy to create complex, multi-step workflows that can be run efficiently.

## Features

- **Directed Acyclic Graphs**: Build workflows as DAGs with nodes and edges
- **Concurrent Execution**: Automatically parallelize independent nodes
- **Data Flow Management**: Pass data seamlessly between workflow steps
- **Cycle Detection**: Prevent infinite loops with built-in cycle detection
- **Configurable Workers**: Control concurrency with customizable worker pools

## Installation

```bash
go get github.com/dskart/gollum/ringchain
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dskart/gollum/ringchain"
	"go.uber.org/zap"
)

// Define custom nodes
type GreetingNode struct{}

func (n *GreetingNode) Name() string {
	return "greeting"
}

func (n *GreetingNode) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	name := args["name"].(string)
	return map[string]any{
		"greeting": fmt.Sprintf("Hello, %s!", name),
	}, nil
}

type EnhancerNode struct{}

func (n *EnhancerNode) Name() string {
	return "enhancer"
}

func (n *EnhancerNode) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	greeting := args["greeting"].(string)
	return map[string]any{
		"enhanced": fmt.Sprintf("%s How are you today?", greeting),
	}, nil
}

func main() {
	// Create a new graph
	graph := ringchain.NewGraph()

	// Add nodes
	graph.AddNode("greeting", &GreetingNode{})
	graph.AddNode("enhancer", &EnhancerNode{})

	// Connect nodes
	graph.AddEdge("greeting", "enhancer")

	// Create a logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Execute the graph
	results, err := graph.Execute(
		context.Background(),
		logger,
		map[string]any{"name": "World"},
		ringchain.WithNumWorkers(2),
	)

	if err != nil {
		fmt.Printf("Error executing graph: %v\n", err)
		os.Exit(1)
	}

	// Print the results from each node
	for nodeHash, result := range results {
		fmt.Printf("Node %s output: %v\n", nodeHash, result)
	}
}
```

## Building Complex Workflows

You can build more complex workflows by adding multiple nodes and connections:

```go
// Create a new graph
graph := ringchain.NewGraph()

// Add nodes
graph.AddNode("data_retrieval", &DataRetrievalNode{})
graph.AddNode("data_processing", &DataProcessingNode{})
graph.AddNode("data_analysis", &DataAnalysisNode{})
graph.AddNode("report_generation", &ReportGenerationNode{})

// Connect nodes
graph.AddEdge("data_retrieval", "data_processing")
graph.AddEdge("data_processing", "data_analysis")
graph.AddEdge("data_analysis", "report_generation")

// Execute the graph
results, err := graph.Execute(
    context.Background(),
    logger,
    initialArgs,
)
```

## Parallel Execution

Ringchain automatically executes independent nodes in parallel:

```go
// Create a new graph with parallel branches
graph := ringchain.NewGraph()

// Add nodes
graph.AddNode("start", &StartNode{})
graph.AddNode("branch_a", &BranchANode{})
graph.AddNode("branch_b", &BranchBNode{})
graph.AddNode("combine", &CombineNode{})

// Connect nodes
graph.AddEdge("start", "branch_a")
graph.AddEdge("start", "branch_b")
graph.AddEdge("branch_a", "combine")
graph.AddEdge("branch_b", "combine")

// Execute with multiple workers
results, err := graph.Execute(
    context.Background(),
    logger,
    initialArgs,
    ringchain.WithNumWorkers(4),
)
```


## Integration with GoLLuM

This package works seamlessly with other GoLLuM modules:

- Use with [OpenAI](../openai) for API communication
- Use with [Scrolls](../scrolls) for prompt templating and management