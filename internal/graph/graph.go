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

// Partition returns a slice of WeightedEdges that should be partitioned (i.e., "broken").
// The pseudo-algorithm for Partition is as follows:
// 1. find root nodes
// 2. find all edges in the graph
// 3. for each root
// 		a. recursively traverse the graph to identify the min shortest path to get to a node
// 		b. for each edge in the graph, store the min shortest path to a node from the root under observation
// 4. calculate edge distances
// 5. identify edges where the sum of cross-iteration distances > epsilon
func (g Graph) Partition(epsilon float64) []WeightedEdge {
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

	dists := g.calculateDistances(edges)
	return g.identifyPartitions(dists, epsilon)
}

// calculateDistances calculates the strongest strong distance of an edge from
// every possible root in the graph.
//
// The strongest strong distance is defined as:
// * 1 / CONN(u, v)
//   * where CONN(u, v) denotes the strength of connectedness of a pair of vertices.
//   * CONN(u, v) is defined as the MAX{s(P)}, where s(P) denotes the MIN edge weight in a path.
//   * CONN(u, v) can be summarized as the max, weakest edge between nodes of various paths.
//
// "Distances in Weighted Graphs" with authors Dhanyamol M V and Sunil Mathew.
//   * http://www.researchmathsci.org/apamart/apam-v8n1-1.pdf
//
// In splitfile, we modify this equation slightly to get a more global understanding
// * 1 / SUM(CONN(u, v))
// 	 * where we are summing the CONN(u, v) from different roots
//
// If you are interested in my notes/highlights, check out this Tweet
//   * https://twitter.com/McCurdyColton/status/1179024173664002049?s=20
func (g Graph) calculateDistances(edges []WeightedEdge) []WeightedEdge {
	res := make([]WeightedEdge, 0, len(edges))

	for _, edge := range edges {

		var sum float64
		for _, mps := range edge.MinPathStrengths {
			sum += mps
		}

		edge.Distance = 1 / (sum)
		res = append(res, edge)
	}

	return res
}

// identifyPartitions identifies the edges where the distance is greater than the
// configured epsilon value. These edges will be split/"broken".
func (g Graph) identifyPartitions(edges []WeightedEdge, epsilon float64) []WeightedEdge {
	res := make([]WeightedEdge, 0, len(edges))

	return res
}

// recursiveBFS does recursive breadth-first search of the graph, keeping track
// of shortest path and the mininimum path strength (i.e., the weaked edge in a path).
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
