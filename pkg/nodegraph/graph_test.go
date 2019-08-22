package nodegraph

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name              string
		expectedRelations map[Node][]Node
	}{
		{
			name:              "empty relations map",
			expectedRelations: map[Node][]Node{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := New()

			assert.Equal(t, test.expectedRelations, actual.relations, "relations map does not match expected.")
		})
	}
}

type mockNode struct {
	pos int
}

func (m mockNode) Pos() token.Pos {
	return token.Pos(m.pos)
}

func TestAddNodes(t *testing.T) {
	tests := []struct {
		name  string
		graph *Graph
		nodes []Node
	}{
		{
			name: "add nil node",
			graph: &Graph{
				relations: make(map[Node][]Node),
			},
			nodes: []Node{nil},
		},
		{
			name: "add single node",
			graph: &Graph{
				relations: make(map[Node][]Node),
			},
			nodes: []Node{mockNode{pos: 1}},
		},
		{
			name: "add multiple nodes",
			graph: &Graph{
				relations: make(map[Node][]Node),
			},
			nodes: []Node{mockNode{pos: 1}, mockNode{pos: 2}, mockNode{pos: 3}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected := make(map[Node][]Node)

			for _, n := range test.nodes {
				expected[n] = make([]Node, 0, 0)
			}

			test.graph.AddNodes(test.nodes...)

			assert.Equal(t, expected, test.graph.relations, "the expected map does not match the graphs relation map.")
		})
	}
}

func TestAddEdges(t *testing.T) {
	tests := []struct {
		name     string
		graph    *Graph
		source   Node
		related  []Node
		expected *Graph
	}{
		{
			name: "nil source",
			graph: &Graph{
				relations: make(map[Node][]Node),
			},
			source:  nil,
			related: []Node{mockNode{pos: 1}},
			expected: &Graph{
				relations: make(map[Node][]Node),
			},
		},
		{
			name: "nil related non-existent source",
			graph: &Graph{
				relations: make(map[Node][]Node),
			},
			source:  mockNode{pos: 1},
			related: nil,
			expected: &Graph{
				relations: map[Node][]Node{mockNode{pos: 1}: make([]Node, 0)},
			},
		},
		{
			name: "nil related source exists",
			graph: &Graph{
				relations: map[Node][]Node{mockNode{pos: 1}: make([]Node, 0)},
			},
			source:  mockNode{pos: 1},
			related: nil,
			expected: &Graph{
				relations: map[Node][]Node{mockNode{pos: 1}: make([]Node, 0)},
			},
		},
		{
			name: "adding single edge to existing node",
			graph: &Graph{
				relations: map[Node][]Node{mockNode{pos: 1}: make([]Node, 0)},
			},
			source:  mockNode{pos: 1},
			related: []Node{mockNode{pos: 2}},
			expected: &Graph{
				relations: map[Node][]Node{
					mockNode{pos: 1}: []Node{mockNode{pos: 2}},
					mockNode{pos: 2}: make([]Node, 0),
				},
			},
		},
		{
			name: "adding single edge to non-existent node",
			graph: &Graph{
				relations: map[Node][]Node{mockNode{pos: 1}: make([]Node, 0)},
			},
			source:  mockNode{pos: 2},
			related: []Node{mockNode{pos: 3}},
			expected: &Graph{
				relations: map[Node][]Node{
					mockNode{pos: 1}: make([]Node, 0),
					mockNode{pos: 2}: []Node{mockNode{pos: 3}},
					mockNode{pos: 3}: make([]Node, 0),
				},
			},
		},
		{
			name: "source and related are the same node",
			graph: &Graph{
				relations: map[Node][]Node{mockNode{pos: 1}: make([]Node, 0)},
			},
			source:  mockNode{pos: 1},
			related: []Node{mockNode{pos: 1}},
			expected: &Graph{
				relations: map[Node][]Node{mockNode{pos: 1}: make([]Node, 0)},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.graph.AddEdges(test.source, test.related...)

			assert.Equal(t, test.expected, test.graph, "the expected graph does not match the actual graph after adding edges.")
		})
	}
}
