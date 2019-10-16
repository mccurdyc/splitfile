package graph

import (
	"errors"
	"go/types"
)

// WeightedEdge contains pointers to source and destination Nodes as well as
// the assigned weight of the relationship.
type WeightedEdge struct {
	Weight float64
	Dest   *Node
}

// Node has an ID and a map of WeightedEdges or weighted relationships to other
// Nodes.
type Node struct {
	ID            string
	Object        types.Object
	Edges         map[string]WeightedEdge
	Paths [][]*Node
	// ConnectednessStrength is the current "strongest" or highest shortest path distance
	ConnectednessStrength float64
}

// NewNode creates a pointer to a new Node with ID, id, and initializes a map of Edges.
func NewNode(id string, obj types.Object) *Node {
	return &Node{
		ID:            id,
		Object:        obj,
		Edges:         make(map[string]WeightedEdge),
		Paths: make([][]*Node, 0), // going to dynamically resize
	}
}

// AddEdges adds weighted edges to a source node that signify relationships with other nodes.
func (n *Node) AddEdge(dest *Node, w float64) {
	// prevent edge between node and itself
	if n.ID == dest.ID {
		return
	}

	n.Edges[dest.ID] = WeightedEdge{
		Weight: w,
		Dest:   dest,
	}
}

// ContainsEdge returns whether or not the graph contains an edge from source to dest.
func (n *Node) ContainsEdge(dest *Node) bool {
	_, ok := n.Edges[dest.ID]
	return ok
}

func (n *Node) Valid() (bool, error) {
	if n == nil {
		return false, errors.New("invalid node; node cannot be nil")
	}

	if len(n.ID) == 0 {
		return false, errors.New("invalid node; node must have an ID")
	}

	if n.Edges == nil {
		return false, errors.New("invalid node; edge map must be initialized")
	}

	return true, nil
}
