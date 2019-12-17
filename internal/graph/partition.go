package graph

import (
	"log"
)

func shouldPartition(g Graph) bool {
	if len(g) <= 1 {
		return false
	}

	// TODO: thought; maybe eventually we store meta data about the graph
	// (e.g., types of edges such as method and if there are zero methods, maybe we consider not partitioning)

	return true
}

// Partition returns a slice of WeightedEdges that should be partitioned (i.e., "broken").
// The pseudo-algorithm for Partition is as follows:
// 1. find root nodes
// 2. recursively traverse the graph
// 		a. keeping track of the current shortest path to a node
//    b. keeping track of all observed shortest paths to a node
// 4. calculate edge distances
// 5. identify edges where the sum of cross-iteration distances > epsilon
func Partition(g Graph, epsilon float64) []WeightedEdge {
	roots := make(map[string]*Node)
	for _, r := range g.Roots() {
		roots[r.ID] = r
	}

	var curr, next map[string]*Node
	curr = roots

	for len(next) > 0 {
		next = bfs(curr)
		curr = next
	}

	dists := calculateDistances(g, g.Edges())
	return identifyPartitions(dists, epsilon)
}

// bfs does breadth-first search of the graph, keeping track of the current and past shortest paths.
func bfs(level map[string]*Node) map[string]*Node {
	next := make(map[string]*Node)

	for _, node := range level {
		for _, childEdge := range node.Edges {
			p, changed := shortestPath(childEdge)

			if changed {
				childEdge.Dest.ShortestPath = p
				childEdge.Dest.ShortestPaths = append(childEdge.Dest.ShortestPaths, p)
				next[childEdge.Dest.ID] = childEdge.Dest
			}
		}
	}

	return next
}

func shortestPath(edge WeightedEdge) (res float64, changed bool) {
	res = edge.Weight

	if edge.Source != nil && edge.Source.ShortestPath != defaultWeight {
		res += edge.Source.ShortestPath
	}

	if edge.Dest == nil || edge.Dest.ShortestPath == defaultWeight || res < edge.Dest.ShortestPath {
		return res, true
	}

	return edge.Dest.ShortestPath, false
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
func calculateDistances(g Graph, edges []WeightedEdge) []WeightedEdge {
	res := make([]WeightedEdge, 0, len(edges))

	for _, edge := range edges {

		var sum float64
		for _, sp := range edge.Dest.ShortestPaths {
			sum += sp
		}

		if sum == 0.0 {
			continue
		}

		edge.Distance = 1 / (sum)
		res = append(res, edge)
	}

	return res
}

// identifyPartitions identifies the edges where the distance is greater than the
// configured epsilon value. These edges will be split/"broken".
func identifyPartitions(edges []WeightedEdge, epsilon float64) []WeightedEdge {
	res := make([]WeightedEdge, 0, len(edges))

	for _, e := range edges {
		if e.Distance > epsilon {
			log.Printf("identified partition with distance: %+v\n\tsrc: %+v \n\tdest: %+v\n", e.Distance, e.Source, e.Dest)
			res = append(res, e)
		}
	}

	return res
}
