package graph

import (
	"errors"
)

const (
	defaultConnectednessValue = -1 // default to an invalid, "infinity", value.
)

// WeightedEdge contains pointers to source and destination Nodes as well as
// the assigned weight of the relationship.
type WeightedEdge struct {
	Weight       float64
	Source, Dest *Node

	ConnectednessStrength float64
	MinPathStrengths      []float64 // used across multiple graph traversals

	Partition bool
	Distance  float64
}

// Node has an ID and a map of WeightedEdges or weighted relationships to other
// Nodes.
type Node struct {
	ID      string
	Object  interface{}
	Edges   map[string]WeightedEdge
	Parents map[string]WeightedEdge // TODO: may be able to delete this now

	MinPathStrength float64
	ShortestPathLen float64 // used for a single graph traversals
}

// NewNode creates a pointer to a new Node with ID, id, and initializes a map of Edges.
func NewNode(id string, v interface{}) *Node {
	return &Node{
		ID:      id,
		Object:  v,
		Edges:   make(map[string]WeightedEdge),
		Parents: make(map[string]WeightedEdge),
	}
}

// AddEdges adds weighted edges to a source node that signify relationships with other nodes.
// Also adds parents to the destination node.
func (n *Node) AddEdge(dest *Node, w float64) {
	// prevent edge between node and itself
	if n.ID == dest.ID {
		return
	}

	n.Edges[dest.ID] = WeightedEdge{
		Weight: w,
		Dest:   dest,
	}

	dest.Parents[n.ID] = WeightedEdge{
		Weight: w,
		Dest:   n,
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

	if n.Parents == nil {
		return false, errors.New("invalid node; parent map must be initialized")
	}

	return true, nil
}
