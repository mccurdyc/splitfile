package splitfile

import (
	"go/ast"

	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/mccurdyc/splitfile/internal/graph"
)

var Analyzer = &analysis.Analyzer{
	Name:     "splitfile",
	Doc:      "checks for clean splits of files in packages based on objects and their relationships with other objects.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	g := traverse(pass.TypesInfo.Defs)

	nodes := g.Partition() // TODO (Issue#8): right now, this returns every single node in the graph
	for _, n := range nodes {
		pass.Reportf(n.Object.Pos(), "parition found - %+v", n)
	}

	return nil, nil
}

// traverse traverses the map of definitions and builds a graph based on the
// relationships.
func traverse(defs map[*ast.Ident]types.Object) graph.Graph {
	g := graph.New()

	for _, def := range defs {

		if skip := filter(def); skip {
			continue
		}

		node := graph.NewNode(def.(types.Object))
		err := g.AddNode(node)
		if err != nil {
			continue
		}

		err = addRelated(g, node)
		if err != nil {
			continue
		}
	}

	return g
}

// filter returns whether or not a def should be filtered out.
func filter(def types.Object) bool {
	if def == nil {
		return true
	}

	return false
}

// addRelated given a graph, g, and root node finds relationships
// with other declarations in the same package and adds them to the graph.
//
// TODO (Issue #15): read value from config or use default
// TODO: check other places for related (e.g., funcs, interfaces, etc.)
func addRelated(g graph.Graph, node *graph.Node) error {
	m := checkMethods(types.NewMethodSet(node.Object.Type()))

	for _, r := range m {
		if r.ID == node.ID {
			continue
		}

		err := g.AddNode(r)
		if err != nil {
			continue
		}

		node.AddEdge(r, 5.0) // TODO (Issue #15): read value from config or use default
	}

	// TODO: check other places for related (e.g., funcs, interfaces, etc.)

	return nil
}

// checkMethods checks methods' signatures for related types.
func checkMethods(mset *types.MethodSet) []*graph.Node {
	rel := make([]*graph.Node, 0)

	for i := 0; i < mset.Len(); i++ {
		method := mset.At(i)

		m := graph.NewNode(method.Obj())
		rel = append(rel, m) // methods themselves are always related

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
func checkSignature(sig *types.Signature) []*graph.Node {
	rel := make([]*graph.Node, 0)

	if v := checkVar(sig.Recv()); v != nil {
		rel = append(rel, v)
	}
	rel = append(rel, checkTuple(sig.Params())...)
	rel = append(rel, checkTuple(sig.Results())...)

	return rel
}

// checkVar validates a variable and if it is valid, it is returned as a valid related.
func checkVar(v *types.Var) *graph.Node {
	if v == nil || v.Type() == types.Type(nil) {
		return nil
	}

	return graph.NewNode(v)
}

// checkTuple checks a tuple of variables for related nodes.
func checkTuple(vars *types.Tuple) []*graph.Node {
	rel := make([]*graph.Node, 0)

	for i := 0; i < vars.Len(); i++ {
		v := checkVar(vars.At(i))

		if v == nil {
			continue
		}

		rel = append(rel, v)
	}

	return rel
}
