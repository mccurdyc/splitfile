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

		related, err := findRelated(node.(types.Object))
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
func findRelated(node types.Object) ([]string, error) {
	rel := make([]string, 0)

	related := checkMethods(types.NewMethodSet(node.Type()))
	rel = append(rel, related...)

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
