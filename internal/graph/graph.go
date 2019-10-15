package graph

import (
	"container/list"

	"github.com/pkg/errors"
)

// Weights defines the configurable edge weights
type Weights struct {
	// TypeField defines the edge weight of a type that has a field of another type.
	TypeField float64
	// TypeMethod defines the edge weight of a type that is a func/method receiver.
	TypeMethod float64
	// VarOfType defines the edge weight between a variable and the variable's type.
	VarOfType float64
	// ConstOfType defines the edge weight between a constant and the constant's type.
	ConstOfType float64
	// FuncParam defines the edge weight between a function parameter and the parameter's type.
	FuncParam float64
	// FuncReturn defines the edge weight between a function return value and the value's type.
	FuncReturn float64
	// FuncBody defines the edge weight between a value used in a function body and the value's type.
	FuncBody float64
}

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
// thoughts:
// 1. we dont necessarily know how the graph is structured; we do know the nodes
//   a. this is why we have to traverse from every node as the root
// 2. we need the distance between every pair of _connected_ nodes
//   a. we have all of the edge weights
//   b. we dont yet know about connectedness (we'll find that out here when we traverse the graph)
// 3. calculate distance of every root -> leaf
//   a. return paths and their weights
func (g Graph) Partition(epsilon float64) [][]Node {
	visited := make(map[string]map[string][]float64)

	for _, root := range g {
		calculateDistance := func(a, b *Node, m map[string]map[string][]float64) {
			m[a.ID][b.ID] = append(m[a.ID][b.ID], .Weight)
		}

		bfs(root, visited)
	}
}

func bfs(root *Node, visited map[string]map[string][]float64, fn func(*Node, *Node, map[string]map[string][]float64)) {
	queue := list.New()
	visited[root.ID] = make(map[string][]float64) // paths [root][edge] and their distances

	for _, edge := range root.Edges {
		if _, ok := visited[edge.Dest.ID]; ok {
			fn(root, edge.Dest, visited)
			continue
		}

		queue.PushBack(edge.Dest)
	}

	for queue.Len() > 0 {
		e := queue.Front()
		n, ok := e.Value.(*Node)
		if !ok {
			continue
		}

		bfs(n, visited, fn)
	}
}

func (g Graph) shouldPartition() bool {
	if len(g) <= 1 {
		return false
	}

	// TODO: thought; maybe eventually we store meta data about the graph
	// (e.g., types of edges such as method and if there are zero methods, maybe we consider not partitioning)

	return true
}
