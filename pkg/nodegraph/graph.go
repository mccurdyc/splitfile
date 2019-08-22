package nodegraph

import (
	"go/token"
	"reflect"
	"sync"
)

type Node interface {
	Pos() token.Pos
}

// Graph represents an undirected graph where nodes are declarations and edges
// signify that one declaration is "related" to another declaration.
//
// Related is defined by the consuming package by creating the edges between nodes.
type Graph struct {
	relations map[Node][]Node
	lock      sync.RWMutex
}

// New creates a pointer to a new Graph where the relationship map is intialized.
func New() *Graph {
	return &Graph{
		relations: make(map[Node][]Node),
	}
}

// AddNodes adds nodes, or declarations, to the Graph.
func (g *Graph) AddNodes(nodes ...Node) {
	g.addNodes(nodes...)
}

// addNodes adds nodes to the Graph with a default relationship slice.
func (g *Graph) addNodes(nodes ...Node) {
	for _, node := range nodes {
		// TODO: should this be on the consumer to check or should we abstract this from them?
		// My only thought is that if we don't protect them, they could overwrite a node
		// that already has edges associated.
		if !g.containsNode(node) {
			g.lock.Lock()
			g.relations[node] = make([]Node, 0, 0)
			g.lock.Unlock()
		}
	}
}

// Nodes returns the slice of added nodes.
func (g *Graph) Nodes() []Node {
	keys := make([]Node, 0, len(g.relations))

	g.lock.RLock()
	defer g.lock.RUnlock()

	for node, _ := range g.relations {
		keys = append(keys, node)
	}

	return keys
}

// ContainsNode returns whether or not the graph contains the provided node.
func (g *Graph) ContainsNode(node Node) bool {
	return g.containsNode(node)
}

// containsNode returns whether or not the graph contains the provided node.
func (g *Graph) containsNode(node Node) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	_, ok := g.relations[node]
	return ok
}

// AddEdges adds edges to a source node that signify relationships with other
// nodes in the Graph.
//
// TODO (@mccurdyc): this could add duplicate edges.
func (g *Graph) AddEdges(source Node, related ...Node) {
	g.addEdges(source, related...)
}

// addEdges adds a single edge to the source node.
func (g *Graph) addEdges(source Node, related ...Node) {
	if source == nil {
		return
	}

	if !g.containsNode(source) {
		g.addNodes(source)
	}

	for _, rel := range related {
		// a node can't be related to itself
		if reflect.DeepEqual(rel, source) {
			continue
		}

		g.addNodes(rel)

		g.lock.Lock()
		g.relations[source] = append(g.relations[source], rel)
		g.lock.Unlock()
	}
}

// EdgesOf returns the nodes of connected edges to a given source node.
func (g *Graph) EdgesOf(source Node) []Node {
	return g.relations[source]
}
