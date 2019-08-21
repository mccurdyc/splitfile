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

func (m *mockNode) Pos() token.Pos {
	return token.Pos(m.pos)
}

func TestAddNodes(t *testing.T) {
	tests := []struct {
		name  string
		graph *Graph
		nodes []Node
	}{
		{
			name: "add single node",
			graph: &Graph{
				relations: make(map[Node][]Node),
			},
			nodes: []Node{&mockNode{pos: 1}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected := make(map[Node][]Node)

			for _, n := range test.nodes {
				expected[n] = make([]Node, 0, 1)
			}

			test.graph.AddNodes(test.nodes...)

			assert.Equal(t, expected, test.graph.relations, "the expected map does not match the graphs relation map.")
		})
	}
}
