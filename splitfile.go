package splitfile

import (
	"errors"
	"fmt"
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

		graph.AddNodes(node)

		// This could be done recursively
		// Right now, the thought was to only add a single level of related nodes.
		related, err := findRelationships(graph, node.(types.Object))
		if err != nil {
			continue
		}

		for _, relNode := range related {
			graph.AddNodes(relNode)
			graph.AddEdges(node, relNode)
		}
	}

	// findSplits()

	return nil, nil
}

// findRelationships given a root declaration, decl, attempts to find relationships
// with other declarations in the same package.
func findRelationships(graph *nodegraph.Graph, node types.Object) ([]types.Object, error) {
	// always check methods of type
	methods := types.NewMethodSet(node.Type())
	_, _ = checkMethodsForRelated(graph, methods)

	return nil, nil
}

func checkMethodsForRelated(graph *nodegraph.Graph, mset *types.MethodSet) ([]types.Object, error) {
	for i := 0; i < mset.Len(); i++ {
		method := mset.At(i)
		_, err := checkMethod(graph, method)
		if err != nil {
			continue
		}

	}

	return nil, nil
}

var errUnprocessableMethodKind = errors.New("cannot process method selection of kind")

func checkMethod(graph *nodegraph.Graph, method *types.Selection) ([]types.Object, error) {
	if method.Kind() != types.MethodVal {
		return nil, errUnprocessableMethodKind
	}

	obj := method.Obj()
	typ := obj.Type()

	fmt.Println(typ)
	fmt.Println(obj)

	// TODO: need to look through params
	// params := typ.Params()

	// TODO: need to look through children scopes for body

	return nil, nil
}
