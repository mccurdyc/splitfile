package splitfile

import (
	"errors"
	"fmt"
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

	switch def.(type) {
	case *types.Func, *types.TypeName:
		return false
	}

	return true
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

	// it is important to add receivers first
	for _, v := range m {
		if !g.ContainsNode(v.method.ID) {
			err := g.AddNode(v.method)
			if err != nil {
				fmt.Println("error adding method")
				continue
			}
		}

		node.AddEdge(g[v.method.ID], 5.0) // TODO (Issue #15): read value from config or use default

		for _, r := range v.other {
			if !g.ContainsNode(r.ID) {
				err := g.AddNode(r)
				if err != nil {
					fmt.Println("error adding other")
					continue
				}
			}

			v.method.AddEdge(r, 2.0) // params, results, etc should be on the method node, not the receiver node
		}
	}

	// TODO: check other places for related (e.g., funcs, interfaces, etc.)

	return nil
}

type methodSetResult map[string]methodResult
type methodResult struct {
	method    *graph.Node
	receivers []*graph.Node
	other     []*graph.Node
}

// checkMethods checks methods' signatures for related types.
func checkMethods(mset *types.MethodSet) methodSetResult {
	res := make(methodSetResult)

	for i := 0; i < mset.Len(); i++ {
		method := mset.At(i)

		sig, ok := method.Type().(*types.Signature)
		if !ok {
			continue // skip evaluating method bodies
		}

		m := graph.NewNode(Id(method.Obj()), method.Obj())
		v := methodResult{
			method:    m,
			receivers: make([]*graph.Node, 0),
			other:     make([]*graph.Node, 0),
		}

		sigRes := checkSignature(sig)
		v.receivers = append(v.receivers, sigRes.receivers...)
		v.other = append(v.other, sigRes.other...)

		res[m.ID] = v
	}

	return res
}

// checkSignature checks a function signature. Specifically, the parameters and the return types.
func checkSignature(sig *types.Signature) methodResult {
	res := methodResult{
		receivers: make([]*graph.Node, 0),
		other:     make([]*graph.Node, 0),
	}

	if v := checkVar(sig.Recv()); v != nil {
		res.receivers = append(res.receivers, v)
	}
	res.other = append(res.other, checkTuple(sig.Params())...)
	res.other = append(res.other, checkTuple(sig.Results())...)

	return res
}

// checkVar validates a variable and if it is valid, it is returned as a valid node.
func checkVar(v *types.Var) *graph.Node {
	if v == nil || v.Type() == types.Type(nil) {
		return nil
	}

	typ, ok := v.Type().(*types.Named)
	if !ok {
		return nil
	}
	obj := typ.Obj()

	return graph.NewNode(Id(obj), obj)
}

// checkTuple checks a tuple of variables for nodes.
func checkTuple(vars *types.Tuple) []*graph.Node {
	res := make([]*graph.Node, 0)

	for i := 0; i < vars.Len(); i++ {
		v := checkVar(vars.At(i))

		if v == nil {
			continue
		}

		res = append(res, v)
	}

	return res
}
