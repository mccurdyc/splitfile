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

	for _, node := range pass.TypesInfo.Defs {
		// nil for a package definition
		if node == nil {
			continue
		}

		nodeKey := node.Type().String()
		graph.AddNodes(nodeKey)

		// This could be done recursively
		// Right now, the thought was to only add a single level of related nodes.
		related, err := findRelationships(graph, node.(types.Object))
		if err != nil {
			continue
		}

		for _, relNode := range related {
			graph.AddNodes(relNode)
			graph.AddEdges(nodeKey, relNode)
		}
	}

	// findSplits()

	return nil, nil
}

// findRelationships given a root declaration, decl, attempts to find relationships
// with other declarations in the same package.
func findRelationships(graph *nodegraph.Graph, node types.Object) ([]string, error) {
	return nil, errors.New("not implemented")
}
