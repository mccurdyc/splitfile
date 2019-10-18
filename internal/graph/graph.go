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
	roots := g.FindRoots()

	edges := make([]WeightedEdge, 0)
	for _, node := range g {
		for _, edge := range node.Edges {
			edges = append(edges, edge)
		}
	}

	for _, root := range roots {
		visited := make(map[string]bool)
		recursiveBFS(root, visited)

		for _, edge := range edges {
			// skip if the node is the root because we don't want to add '0' to the slice of
			// seen min paths. This could lead to division by zero.
			if edge.Source == root {
				continue
			}
			edge.MinPathStrengths = append(edge.MinPathStrengths, edge.Source.MinPathStrength)
		}
	}

	g.calculateDistances(edges)
	g.identifyPartitions(edges, epsilon)
}

func (g Graph) identifyPartitions(epsilon float64) {
}

func (g Graph) calculateDistances() {
	for _, node := range g {
		for _, edge := range node.Edges {

			var sum float64
			for _, mps := range edge.MinPathStrengths {
				sum += mps
			}

			edge.Distance = 1 / (sum)
		}
	}
}

func recursiveBFS(node *Node, visited map[string]bool) {
	queue := list.New()
	visited[node.ID] = true

	for _, edge := range node.Edges {
		if _, ok := visited[edge.Dest.ID]; ok {
			var updateChildren bool

			if edge.Weight < edge.Dest.MinPathStrength {
				edge.Dest.MinPathStrength = edge.Weight
				updateChildren = true
			}

			if node.MinPathStrength < edge.Dest.MinPathStrength {
				edge.Dest.MinPathStrength = node.ShortestPathLen
				updateChildren = true
			}

			pLen := node.ShortestPathLen + edge.Weight
			if pLen < edge.Dest.ShortestPathLen {
				edge.Dest.ShortestPathLen = pLen
				updateChildren = true
			}

			if updateChildren {
				queue.PushBack(edge.Dest)
			}

			continue
		}

		queue.PushBack(edge.Dest)
	}

	for queue.Len() > 0 {
		n := queue.Front()
		e, ok := n.Value.(*Node)
		if !ok {
			continue
		}

		e.ShortestPathLen = node.ShortestPathLen + node.Edges[e.ID].Weight
		e.MinPathStrength = node.Edges[e.ID].Weight

		recursiveBFS(e, visited)
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
