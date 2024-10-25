package ringchain

import (
	"fmt"
	"sync"
)

type Store interface {
	// AddNode should add the given node with the given hash value and node properties to the
	// graph. If the node already exists, it is up to you whether ErrNodeAlreadyExists or no
	// error should be returned.
	AddNode(hash string, value Node) error

	// Node should return the node with the given hash value. If the
	// node doesn't exist, ErrNodeNotFound should be returned.
	Node(hash string) (Node, error)

	// RemoveNode should remove the node with the given hash value. If the node doesn't
	// exist, ErrNodeNotFound should be returned. If the node has edges to other vertices,
	// ErrNodeHasEdges should be returned.
	RemoveNode(hash string) error

	// ListNodes should return all vertices in the graph in a slice.
	ListNodes() ([]string, error)

	// NodeCount should return the number of vertices in the graph. This should be equal to the
	// length of the slice returned by ListNodes.
	NodeCount() (int, error)

	// AddEdge should add an edge between the vertices with the given source and target hashes.
	//
	// If either node doesn't exit, ErrNodeNotFound should be returned for the respective
	// node. If the edge already exists, ErrEdgeAlreadyExists should be returned.
	AddEdge(sourceHash, targetHash string, edge Edge) error

	// UpdateEdge should update the edge between the given vertices with the data of the given
	// Edge instance. If the edge doesn't exist, ErrEdgeNotFound should be returned.
	UpdateEdge(sourceHash string, targetHash string, edge Edge) error

	// RemoveEdge should remove the edge between the vertices with the given source and target
	// hashes.
	//
	// If either node doesn't exist, it is up to you whether ErrNodeNotFound or no error should
	// be returned. If the edge doesn't exist, it is up to you whether ErrEdgeNotFound or no error
	// should be returned.
	RemoveEdge(sourceHash string, targetHash string) error

	// Edge should return the edge joining the vertices with the given hash values. It should
	// exclusively look for an edge between the source and the target node, not vice versa. The
	// graph implementation does this for undirected graphs itself.
	//
	// Note that unlike Graph.Edge, this function is supposed to return an Edge[string], i.e. an edge
	// that only contains the node hashes instead of the vertices themselves.
	//
	// If the edge doesn't exist, ErrEdgeNotFound should be returned.
	Edge(sourceHash string, targetHash string) (Edge, error)

	// ListEdges should return all edges in the graph in a slice.
	ListEdges() ([]Edge, error)

	CreatesCycle(source, target string) (bool, error)
}

type memoryStore struct {
	lock     sync.RWMutex
	vertices map[string]Node

	// outEdges and inEdges store all outgoing and ingoing edges for all vertices. For O(1) access,
	// these edges themselves are stored in maps whose keys are the hashes of the target vertices.
	outEdges map[string]map[string]Edge // source -> target
	inEdges  map[string]map[string]Edge // target -> source
}

func newMemoryStore() Store {
	return &memoryStore{
		vertices: make(map[string]Node),
		outEdges: make(map[string]map[string]Edge),
		inEdges:  make(map[string]map[string]Edge),
	}
}

func (s *memoryStore) AddNode(k string, t Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.vertices[k]; ok {
		return ErrNodeAlreadyExists
	}

	s.vertices[k] = t

	return nil
}

func (s *memoryStore) ListNodes() ([]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	hashes := make([]string, 0, len(s.vertices))
	for k := range s.vertices {
		hashes = append(hashes, k)
	}

	return hashes, nil
}

func (s *memoryStore) NodeCount() (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.vertices), nil
}

func (s *memoryStore) Node(k string) (Node, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.vertices[k]
	if !ok {
		return v, ErrNodeNotFound
	}

	return v, nil
}

func (s *memoryStore) RemoveNode(k string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if _, ok := s.vertices[k]; !ok {
		return ErrNodeNotFound
	}

	if edges, ok := s.inEdges[k]; ok {
		if len(edges) > 0 {
			return ErrNodeHasEdges
		}
		delete(s.inEdges, k)
	}

	if edges, ok := s.outEdges[k]; ok {
		if len(edges) > 0 {
			return ErrNodeHasEdges
		}
		delete(s.outEdges, k)
	}

	delete(s.vertices, k)
	return nil
}

func (s *memoryStore) AddEdge(sourceHash, targetHash string, edge Edge) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.outEdges[sourceHash]; !ok {
		s.outEdges[sourceHash] = make(map[string]Edge)
	}

	s.outEdges[sourceHash][targetHash] = edge

	if _, ok := s.inEdges[targetHash]; !ok {
		s.inEdges[targetHash] = make(map[string]Edge)
	}

	s.inEdges[targetHash][sourceHash] = edge

	return nil
}

func (s *memoryStore) UpdateEdge(sourceHash string, targetHash string, edge Edge) error {
	if _, err := s.Edge(sourceHash, targetHash); err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.outEdges[sourceHash][targetHash] = edge
	s.inEdges[targetHash][sourceHash] = edge

	return nil
}

func (s *memoryStore) RemoveEdge(sourceHash, targetHash string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.inEdges[targetHash], sourceHash)
	delete(s.outEdges[sourceHash], targetHash)
	return nil
}

func (s *memoryStore) Edge(sourceHash, targetHash string) (Edge, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	sourceEdges, ok := s.outEdges[sourceHash]
	if !ok {
		return Edge{}, ErrEdgeNotFound
	}

	edge, ok := sourceEdges[targetHash]
	if !ok {
		return Edge{}, ErrEdgeNotFound
	}

	return edge, nil
}

func (s *memoryStore) ListEdges() ([]Edge, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	res := make([]Edge, 0)
	for _, edges := range s.outEdges {
		for _, edge := range edges {
			res = append(res, edge)
		}
	}
	return res, nil
}

func (s *memoryStore) CreatesCycle(source, target string) (bool, error) {
	if _, err := s.Node(source); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", source, err)
	}

	if _, err := s.Node(target); err != nil {
		return false, fmt.Errorf("could not get vertex with hash %v: %w", target, err)
	}

	if source == target {
		return true, nil
	}

	stack := newStack[string]()
	visited := make(map[string]struct{})

	stack.push(source)

	for !stack.isEmpty() {
		currentHash, _ := stack.pop()

		if _, ok := visited[currentHash]; !ok {
			// If the adjacent vertex also is the target vertex, the target is a
			// parent of the source vertex. An edge would introduce a cycle.
			if currentHash == target {
				return true, nil
			}

			visited[currentHash] = struct{}{}

			for adjacency := range s.inEdges[currentHash] {
				stack.push(adjacency)
			}
		}
	}

	return false, nil
}
