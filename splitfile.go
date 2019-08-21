package splitfile

import (
	"errors"
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

	for _, v := range pass.TypesInfo.Defs {
		// v is nil for a package definition
		if v == nil {
			continue
		}

		graph.AddNodes(v)
	}

	for _, node := range graph.Nodes() {
		err := findRelationships(graph, node.(types.Object))
		if err != nil {
			continue
		}
	}

	return nil, nil
}

// findRelationships given a root declaration, decl, attempts to find relationships
// with other declarations in the same package.
func findRelationships(graph *nodegraph.Graph, node types.Object) error {
	return errors.New("not implemented")
}
