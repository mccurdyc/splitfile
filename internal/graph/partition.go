package graph

import (
	"container/list"
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
// 2. find all edges in the graph
// 3. for each root
// 		a. recursively traverse the graph to identify the min shortest path to get to a node
// 		b. for each edge in the graph, store the min shortest path to a node from the root under observation
// 4. calculate edge distances
// 5. identify edges where the sum of cross-iteration distances > epsilon
func Partition(g Graph, epsilon float64) []WeightedEdge {
	roots := make(map[string]*Node)
	for _, r := range g.Roots() {
		roots[r.ID] = r
	}

	visited := make(map[string]bool)
	recursiveBFS(roots, visited)

	edges := g.Edges()
	for _, edge := range edges {
		if edge.Source.MinPathStrength == defaultWeight {
			continue
		}

		mps := g[edge.Source.ID].Edges[edge.Dest.ID].MinPathStrengths
		mps = append(mps, edge.Source.MinPathStrength)
	}

	dists := calculateDistances(g, edges)
	return identifyPartitions(dists, epsilon)
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
		for _, mps := range edge.MinPathStrengths {
			sum += mps
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

// recursiveBFS does recursive breadth-first search of the graph, keeping track
// of shortest path and the mininimum path strength (i.e., the weaked edge in a path).
func recursiveBFS(level map[string]*Node, visited map[string]bool) {
	evalQueue := list.New()
	nextLevel := make(map[string]*Node)

	for _, node := range level {
		if _, ok := visited[node.ID]; !ok {
			evalQueue.PushBack(node)
			continue
		}

		var reevaluate bool
		for _, parentEdge := range node.Parents {
			if parentEdge.Weight < node.MinPathStrength {
				node.MinPathStrength = parentEdge.Weight
				reevaluate = true
			}

			parent := parentEdge.Source

			if parent.MinPathStrength < node.MinPathStrength {
				node.MinPathStrength = parent.MinPathStrength
				reevaluate = true
			}

			pathLen := parent.ShortestPathLen + parentEdge.Weight
			if pathLen < node.ShortestPathLen {
				node.ShortestPathLen = pathLen
				reevaluate = true
			}

			if reevaluate {
				evalQueue.PushBack(node)
			}
		}
	}

	for evalQueue.Len() > 0 {
		n, ok := evalQueue.Remove(evalQueue.Front()).(*Node)
		if !ok {
			continue
		}

		for _, childEdge := range n.Edges {
			nextLevel[childEdge.Dest.ID] = childEdge.Dest
		}

		for _, parentEdge := range n.Parents {
			if n.MinPathStrength == defaultWeight || parentEdge.Weight < n.MinPathStrength {
				n.MinPathStrength = parentEdge.Weight
			}

			parent := parentEdge.Source
			if parent.MinPathStrength < n.MinPathStrength {
				n.MinPathStrength = parent.MinPathStrength
			}

			pathLen := parent.ShortestPathLen + parentEdge.Weight
			if n.ShortestPathLen == defaultWeight || pathLen < n.ShortestPathLen {
				n.ShortestPathLen = pathLen
			}
		}

		visited[n.ID] = true
		recursiveBFS(nextLevel, visited)
	}
}
