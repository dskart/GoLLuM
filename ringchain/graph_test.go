package ringchain

import (
	"context"
	"fmt"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type TestNode struct {
	name string
}

func (n TestNode) Name() string {
	return n.name
}

func (n TestNode) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	result := maps.Clone(args)

	key := "n_" + n.name
	result[key] = true
	return result, nil
}

func TestGraph_AddNode(t *testing.T) {
	t.Run("AddNodes", func(t *testing.T) {
		nodes := []Node{TestNode{name: "1"}, TestNode{name: "2"}}
		expectedNodeHashes := []string{"1", "2"}
		graph := NewGraph()
		for _, node := range nodes {
			err := graph.AddNode(node.Name(), node)
			require.NoError(t, err)

			retNode, err := graph.Node(node.Name())
			require.NoError(t, err)
			assert.Equal(t, node, retNode)
		}

		store := graph.store
		nodeHashes, err := store.ListNodes()
		require.NoError(t, err)
		assert.ElementsMatch(t, expectedNodeHashes, nodeHashes)
	})

	t.Run("DuplicateNode", func(t *testing.T) {
		graph := NewGraph()

		err := graph.AddNode("1", TestNode{name: "1"})
		require.NoError(t, err)
		err = graph.AddNode("2", TestNode{name: "2"})
		require.NoError(t, err)

		err = graph.AddNode("2", TestNode{name: "2"})
		require.ErrorIs(t, err, ErrNodeAlreadyExists)

		expectedNodeHashes := []string{"1", "2"}
		store := graph.store
		nodeHashes, err := store.ListNodes()
		require.NoError(t, err)
		assert.ElementsMatch(t, expectedNodeHashes, nodeHashes)
	})
}

func TestGraph_Execute(t *testing.T) {
	logger := zaptest.NewLogger(t)

	t.Run("Simple", func(t *testing.T) {
		ctx := context.Background()

		g := NewGraph()
		err := g.AddNode("1", TestNode{name: "1"})
		require.NoError(t, err)
		err = g.AddNode("2", TestNode{name: "2"})
		require.NoError(t, err)

		err = g.AddEdge("1", "2")
		require.NoError(t, err)

		require.NoError(t, err)

		res, err := g.Execute(ctx, logger, map[string]any{"path": []string{}})
		require.NoError(t, err)
		lastNodeRes, ok := res["2"]
		require.True(t, ok)
		assert.Contains(t, lastNodeRes, "n_1")
		assert.Contains(t, lastNodeRes, "n_2")
	})

	t.Run("Parallel", func(t *testing.T) {
		ctx := context.Background()

		g := NewGraph()
		err := g.AddNode("1", TestNode{name: "1"})
		require.NoError(t, err)
		err = g.AddNode("2", TestNode{name: "2"})
		require.NoError(t, err)
		err = g.AddNode("3", TestNode{name: "3"})
		require.NoError(t, err)
		err = g.AddNode("4", TestNode{name: "4"})
		require.NoError(t, err)

		err = g.AddEdge("1", "2")
		require.NoError(t, err)
		err = g.AddEdge("1", "3")
		require.NoError(t, err)

		err = g.AddEdge("2", "4")
		require.NoError(t, err)
		err = g.AddEdge("3", "4")
		require.NoError(t, err)

		require.NoError(t, err)

		res, err := g.Execute(ctx, logger, map[string]any{})
		require.NoError(t, err)
		lastNodeRes, ok := res["4"]
		require.True(t, ok)
		assert.Contains(t, lastNodeRes, "n_1")
		assert.Contains(t, lastNodeRes, "n_2")
		assert.Contains(t, lastNodeRes, "n_3")
		assert.Contains(t, lastNodeRes, "n_4")
		assert.Len(t, lastNodeRes, 4)
	})
}

type ErrorNode struct {
	name string
}

func (n ErrorNode) Name() string {
	return n.name
}

func (n ErrorNode) Run(ctx context.Context, logger *zap.Logger, args map[string]any) (map[string]any, error) {
	return nil, fmt.Errorf("node is broken.")
}

func TestGraphNodeError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	g := NewGraph()
	err := g.AddNode("1", ErrorNode{name: "1"})
	require.NoError(t, err)
	err = g.AddNode("2", ErrorNode{name: "2"})
	require.NoError(t, err)
	err = g.AddNode("3", ErrorNode{name: "3"})
	require.NoError(t, err)
	err = g.AddNode("4", ErrorNode{name: "4"})
	require.NoError(t, err)

	err = g.AddEdge("1", "2")
	require.NoError(t, err)
	err = g.AddEdge("1", "3")
	require.NoError(t, err)

	err = g.AddEdge("2", "4")
	require.NoError(t, err)
	err = g.AddEdge("3", "4")
	require.NoError(t, err)

	require.NoError(t, err)

	_, err = g.Execute(ctx, logger, map[string]any{})
	assert.Equal(t, err.Error(), "node is broken.")
}
