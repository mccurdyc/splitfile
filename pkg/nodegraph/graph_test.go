package nodegraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name              string
		expectedRelations map[string][]string
	}{
		{
			name:              "empty relations map",
			expectedRelations: map[string][]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := New()

			assert.Equal(t, test.expectedRelations, actual.relations, "relations map does not match expected.")
		})
	}
}

func TestAddNodes(t *testing.T) {
	tests := []struct {
		name  string
		graph *Graph
		nodes []string
	}{
		{
			name: "add nil node",
			graph: &Graph{
				relations: make(map[string][]string),
			},
			nodes: []string{""},
		},
		{
			name: "add single node",
			graph: &Graph{
				relations: make(map[string][]string),
			},
			nodes: []string{"a"},
		},
		{
			name: "add multiple nodes",
			graph: &Graph{
				relations: make(map[string][]string),
			},
			nodes: []string{"a", "b", "c"},
		},
		{
			name: "add duplicate nodes",
			graph: &Graph{
				relations: map[string][]string{"a": make([]string, 0)},
			},
			nodes: []string{"a", "a", "a"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected := make(map[string][]string)

			for _, n := range test.nodes {
				expected[n] = make([]string, 0)
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
		source   string
		related  []string
		expected *Graph
	}{
		{
			name: "nil source",
			graph: &Graph{
				relations: make(map[string][]string),
			},
			source:  "",
			related: []string{"a"},
			expected: &Graph{
				relations: make(map[string][]string),
			},
		},
		{
			name: "nil related non-existent source",
			graph: &Graph{
				relations: make(map[string][]string),
			},
			source:  "a",
			related: []string{""},
			expected: &Graph{
				relations: map[string][]string{"a": make([]string, 0)},
			},
		},
		{
			name: "nil related source exists",
			graph: &Graph{
				relations: map[string][]string{"a": make([]string, 0)},
			},
			source:  "a",
			related: []string{""},
			expected: &Graph{
				relations: map[string][]string{"a": make([]string, 0)},
			},
		},
		{
			name: "adding single edge to existing node",
			graph: &Graph{
				relations: map[string][]string{"a": make([]string, 0)},
			},
			source:  "a",
			related: []string{"b"},
			expected: &Graph{
				relations: map[string][]string{
					"a": []string{"b"},
					"b": []string{"a"},
				},
			},
		},
		{
			name: "adding single edge to non-existent node with non-empty graph",
			graph: &Graph{
				relations: map[string][]string{"a": make([]string, 0)},
			},
			source:  "b",
			related: []string{"c"},
			expected: &Graph{
				relations: map[string][]string{
					"a": []string{},
					"b": []string{"c"},
					"c": []string{"b"},
				},
			},
		},
		{
			name: "source and related are the same node",
			graph: &Graph{
				relations: map[string][]string{"a": make([]string, 0)},
			},
			source:  "a",
			related: []string{"a"},
			expected: &Graph{
				relations: map[string][]string{"a": []string{}},
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
