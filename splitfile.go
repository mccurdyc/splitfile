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

		if !graph.ContainsNode(nodeKey) {
			graph.AddNodes(nodeKey)
		}

		related, err := findRelated(graph, node.(types.Object))
		if err != nil {
			continue
		}

		for _, rel := range related {
			if rel == "" || rel == nodeKey {
				continue
			}

			if !graph.ContainsNode(rel) {
				graph.AddNodes(rel)
			}

			graph.AddEdges(nodeKey, rel)
		}
	}

	// findSplits()

	return nil, nil // TODO: FIX THIS!
}

// findRelated given a root node attempts to find relationships
// with other declarations in the same package.
func findRelated(graph *nodegraph.Graph, node types.Object) ([]string, error) {
	rel := make([]string, 0)

	related := checkMethods(graph, types.NewMethodSet(node.Type()))
	rel = append(rel, related...)

	// TODO: check other places for related (e.g., funcs, interfaces, etc.)

	return rel, nil
}

// checkMethods checks methods' signatures for related types.
func checkMethods(graph *nodegraph.Graph, mset *types.MethodSet) []string {
	rel := make([]string, 0)

	for i := 0; i < mset.Len(); i++ {
		method := mset.At(i)
		rel = append(rel, method.String()) // methods themselves are always related

		sig, ok := method.Type().(*types.Signature)
		if !ok {
			continue
		}

		related := checkSignature(graph, sig)
		rel = append(rel, related...)
	}

	return rel
}

// checkSignature checks a function signature, the receiver (if it is a method this
// will be a non-nil value), the parameters and the return types.
func checkSignature(graph *nodegraph.Graph, sig *types.Signature) []string {
	rel := make([]string, 0)

	rel = append(rel, checkVar(graph, sig.Recv()))
	rel = append(rel, checkTuple(graph, sig.Params())...)
	rel = append(rel, checkTuple(graph, sig.Results())...)

	return rel
}

// checkVar checks a variable to see if it is contained in the graph.
func checkVar(graph *nodegraph.Graph, v *types.Var) string {
	if v == nil {
		return ""
	}

	return v.Type().String()
}

// checkTuple checks a tuple of variables to see if they are contained in the graph.
func checkTuple(graph *nodegraph.Graph, vars *types.Tuple) []string {
	rel := make([]string, 0)

	for i := 0; i < vars.Len(); i++ {
		v := checkVar(graph, vars.At(i))

		if v == "" {
			continue
		}

		rel = append(rel, v)
	}

	return rel
}
