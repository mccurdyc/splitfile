package splitfile

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/mccurdyc/splitfile/internal/graph"
)

var Analyzer = &analysis.Analyzer{
	Name:     "splitfile",
	Doc:      "A static analysis that identifies partitions of declarations and their uses to improve the readability of Go packages.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

type Poser interface {
	Pos() token.Pos
}

func run(pass *analysis.Pass) (interface{}, error) {
	g := traverse(pass.TypesInfo.Defs)

	edges := graph.Partition(g, 0.1) // TODO: make this epsilon value configurable (or at least find a reasonable default i.e., the "natural" value)
	for _, e := range edges {
		src, ok := e.Source.Object.(Poser)
		if !ok {
			continue
		}

		dest, ok := e.Dest.Object.(Poser)
		if !ok {
			continue
		}

		pass.Reportf(src.Pos(), "parition found between -> %+v", dest.Pos())
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

		// if the graph already contains a node with the same ID, no need to create
		// another node with the same ID. This is because multiple variables in the
		// same package could use the same name and have the same type. We don't want
		// to skip entirely because this is a new instance of that variable and should
		// be checked for new related nodes.
		if !g.ContainsNode(Id(def)) {
			node := graph.NewNode(Id(def), def)
			err := g.AddNode(node)
			if err != nil {
				continue
			}
		}

		err := addRelated(g, g[Id(def)])
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

type Typer interface {
	Type() types.Type
}

// addRelated given a graph, g, and root node finds relationships
// with other declarations in the same package and adds them to the graph.
//
// TODO (Issue #15): read value from config or use default
// TODO: check other places for related (e.g., funcs, interfaces, etc.)
func addRelated(g graph.Graph, node *graph.Node) error {
	t, ok := node.Object.(Typer)
	if !ok {
		return errors.New("node does not have a Type() method")
	}
	m := checkMethods(types.NewMethodSet(t.Type()))

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

		m := graph.NewNode(Id(method.Obj()), method.Obj())
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

// checkSignature checks a function signature. Specifically, the parameters and the return types.
func checkSignature(sig *types.Signature) []*graph.Node {
	rel := make([]*graph.Node, 0)

	rel = append(rel, checkTuple(sig.Params())...)
	rel = append(rel, checkTuple(sig.Results())...)

	return rel
}

// checkVar validates a variable and if it is valid, it is returned as a valid related.
func checkVar(v *types.Var) *graph.Node {
	if v == nil || v.Type() == types.Type(nil) {
		return nil
	}

	return graph.NewNode(Id(v), v)
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
