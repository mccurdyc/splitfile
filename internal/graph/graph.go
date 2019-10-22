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

// FindRoots finds the many possible roots in the graph.
func (g Graph) FindRoots() []*Node {
	roots := make([]*Node, 0, len(g))

	for _, node := range g {
		if len(node.Parents) == 0 {
			roots = append(roots, node)
		}
	}

	return roots
}
