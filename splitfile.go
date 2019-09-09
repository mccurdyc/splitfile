package splitfile

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/mccurdyc/splitfile/pkg/nodegraph"
)

var Analyzer = &analysis.Analyzer{
	Name:     "splitfile",
	Doc:      "checks for clean splits of files in packages based on objects and their relationships with other objects.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	graph := nodegraph.New()

	for _, node := range pass.TypesInfo.Defs {
		// nil for a package definition
		if node == nil {
			continue
		}

		nodeKey := node.Type().String()
		graph.AddNodes(nodeKey)

		related, err := findRelated(node.(types.Object))
		if err != nil {
			continue
		}

		for _, rel := range related {
			if rel == "" || rel == nodeKey {
				continue
			}

			graph.AddNodes(rel)
			graph.AddEdges(nodeKey, rel)
		}
	}

	findPartitions(graph)

	return nil, nil
}

// findPartitions traverses the graph --- using breadth-first search --- checking
// for partitions using the geodesic (shortest-path) edge betweeness divisive algorithm of Girvan and Newman.
// Since we are analyzing an "unweighted" graph --- i.e., a graph where all edges are equal --- we
// assign equal weights of 1 to all edges for the shortest-path calculation.
//
// As described in the Fortunato paper, the algorithm behaves as follows:
//
//		1. Computation of the centrality for all edges
// 		2. Removal of edge with largest centrality: in case of ties with other edges, one of them is picked at random
// 		3. Recalculation of centralities on the running graph
// 		4. Iteration of the cycle from step 2
//
// Community detection in graphs - Santo Fortunato
// 	* https://arxiv.org/abs/0906.0612
//
// Betweenness-based decomposition methods for social and biological networks
// 	(Uses edge and vertex betweeness) for overlapping communities
// 	* http://www1.maths.leeds.ac.uk/Statistics/workshop/lasr2006/proceedings/pinney-talk.pdf
//
// A Faster Algorithm for Betweenness Centrality - https://kops.uni-konstanz.de/bitstream/handle/123456789/5739/algorithm.pdf
// Brandes' algorithm:
// 	* https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.72.9610&rep=rep1&type=pdf
// 	* (unofficial paper) https://www.cl.cam.ac.uk/teaching/1617/MLRD/handbook/brandes.pdf
func findPartitions(graph *nodegraph.Graph) {
	// iterate through keys (nodes)
	// need to find uniqueness in the related edges for a given node
}

// findRelated given a root node attempts to find relationships
// with other declarations in the same package.
func findRelated(node types.Object) ([]string, error) {
	rel := make([]string, 0)

	related := checkMethods(types.NewMethodSet(node.Type()))
	for _, r := range related {
		if r == node.Type().String() {
			continue
		}
		rel = append(rel, r)
	}

	// TODO: check other places for related (e.g., funcs, interfaces, etc.)

	return rel, nil
}

// checkMethods checks methods' signatures for related types.
func checkMethods(mset *types.MethodSet) []string {
	rel := make([]string, 0)

	for i := 0; i < mset.Len(); i++ {
		method := mset.At(i)
		rel = append(rel, method.String()) // methods themselves are always related

		sig, ok := method.Type().(*types.Signature)
		if !ok {
			continue
		}

		related := checkSignature(sig)
		rel = append(rel, related...)
	}

	return rel
}

// checkSignature checks a function signature, the receiver (if it is a method this
// will be a non-nil value), the parameters and the return types.
func checkSignature(sig *types.Signature) []string {
	rel := make([]string, 0)

	if v := checkVar(sig.Recv()); v != "" {
		rel = append(rel, v)
	}
	rel = append(rel, checkTuple(sig.Params())...)
	rel = append(rel, checkTuple(sig.Results())...)

	return rel
}

// checkVar validates a variable and if it is valid, it is returned as a valid related.
func checkVar(v *types.Var) string {
	if v == nil || v.Type() == types.Type(nil) {
		return ""
	}

	var res string

	switch t := v.Type().(type) {
	case *types.Slice:
		res = t.Elem().String()
	default:
		res = t.String()
	}

	return res
}

// checkTuple checks a tuple of variables for related nodes.
func checkTuple(vars *types.Tuple) []string {
	rel := make([]string, 0)

	for i := 0; i < vars.Len(); i++ {
		v := checkVar(vars.At(i))

		if v == "" {
			continue
		}

		rel = append(rel, v)
	}

	return rel
}
