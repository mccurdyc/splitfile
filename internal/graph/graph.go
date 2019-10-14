package graph

import (
	"github.com/pkg/errors"
)

// Graph is a map of node IDs to the Node with that ID.
type Graph map[string]*Node

// New creates a pointer to a Graph and initializes a map of nodes.
func New() Graph {
	g := make(map[string]*Node)
	return Graph(g)
}

// AddNode adds a valid node to the Graph.
func (g Graph) AddNode(n *Node) error {
	if ok, err := n.Valid(); !ok {
		return errors.Wrap(err, "could not add invalid node")
	}

	g[n.ID] = n
	return nil
}

// ContainsNode returns whether or not the graph contains a node with the given id.
func (g Graph) ContainsNode(id string) bool {
	_, ok := g[id]
	return ok
}

// Partition returns a slice of nodes that should be split from a given source graph.
//
// TODO (Issue#8): actually add logic to properly partition
func (g Graph) Partition() []*Node {
	res := make([]*Node, 0, len(g))

	// TODO: for now, just return every node
	for _, v := range g {
		res = append(res, v)
	}

	return res
}

func (g Graph) shouldPartition() bool {
	if len(g) <= 1 {
		return false
	}

	// TODO: thought; maybe eventually we store meta data about the graph
	// (e.g., types of edges such as method and if there are zero methods, maybe we consider not partitioning)

	return true
}
