package objectgraph

import (
	"go/types"
	"sync"
)

// Graph represents an undirected graph where nodes are declarations and edges
// signify that one declaration is "related" to another declaration.
//
// Related is defined by the consuming package by creating the edges between nodes.
type Graph struct {
	relations map[types.Object][]types.Object
	lock      sync.RWMutex
}

// New creates a pointer to a new Graph where the relationship map is intialized.
func New() *Graph {
	return &Graph{
		relations: make(map[types.Object][]types.Object),
	}
}

// AddNodes adds nodes, or declarations, to the Graph.
func (g *Graph) AddNodes(nodes ...types.Object) {
	g.addNodes(nodes...)
}

// addNodes adds nodes, or declarations, to the Graph.
//
// TODO (@mccurdyc): this could add duplicate nodes.
func (g *Graph) addNodes(nodes ...types.Object) {
	for _, node := range nodes {
		g.addEdges(node, make([]types.Object, 0, 1)...)
	}
}

// Nodes returns the slice of added nodes.
func (g *Graph) Nodes() []types.Object {
	keys := make([]types.Object, 0, len(g.relations))

	for node, _ := range g.relations {
		keys = append(keys, node)
	}

	return keys
}

// ContainsNode returns whether or not the graph contains the provided node.
func (g *Graph) ContainsNode(node types.Object) bool {
	_, ok := g.relations[node]
	return ok
}

// AddEdges adds edges to a source node that signify relationships with other
// nodes in the Graph.
//
// TODO (@mccurdyc): this could add duplicate edges.
func (g *Graph) AddEdges(source types.Object, related ...types.Object) {
	g.addEdges(source, related...)
}

// addEdges adds a single edge to the source node.
func (g *Graph) addEdges(source types.Object, related ...types.Object) {
	g.lock.Lock()
	g.relations[source] = append(g.relations[source], related...)
	g.lock.Unlock()
}

// EdgesOf returns the nodes of connected edges to a given source node.
func (g *Graph) EdgesOf(source types.Object) []types.Object {
	return g.relations[source]
}
