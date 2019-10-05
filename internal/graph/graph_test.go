package graph

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Graph
	}{
		{
			name: "empty relations map",
			want: Graph(map[string]*Node{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New()

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() = %+v mismatch:\n%+v", got, test.want)
			}
		})
	}
}

var (
	edgeB = map[string]WeightedEdge{
		"b": WeightedEdge{
			Weight: 1.0,
			Dest:   &Node{ID: "b", Edges: make(map[string]WeightedEdge)},
		},
	}

	nodeAWithSingleEdgeB = Node{ID: "a", Edges: edgeB}
	nodeCWithSingleEdgeB = Node{ID: "c", Edges: edgeB}
)

func TestAddNode(t *testing.T) {
	tests := []struct {
		name    string
		graph   Graph
		node    *Node
		want    Graph
		wantErr error
	}{
		{
			name:    "add-nil-node-to-empty-graph-should-return-empty-graph",
			graph:   Graph(map[string]*Node{}),
			node:    &Node{},
			want:    Graph(map[string]*Node{}),
			wantErr: errors.New("could not add invalid node: invalid node; node cannot be nil"),
		},

		{
			name:    "add-node-with-nil-ID-to-empty-graph-should-return-empty-graph",
			graph:   Graph(map[string]*Node{}),
			node:    &Node{ID: ""},
			want:    Graph(map[string]*Node{}),
			wantErr: errors.New("could not add invalid node: invalid node; node must have an ID"),
		},

		{
			name:    "add-node-with-nil-edge-map-to-empty-graph-should-return-empty-graph",
			graph:   Graph(map[string]*Node{}),
			node:    &Node{ID: "abc"},
			want:    Graph(map[string]*Node{}),
			wantErr: errors.New("could not add invalid node: invalid node; edge map must be initializated"),
		},

		{
			name:    "add-valid-node-to-empty-graph-should-return-graph-with-that-node",
			graph:   Graph(map[string]*Node{}),
			node:    &nodeAWithSingleEdgeB,
			want:    Graph(map[string]*Node{"a": &nodeAWithSingleEdgeB}),
			wantErr: nil,
		},

		{
			name:    "add-duplicate-node-to-graph-should-return-same-graph",
			graph:   Graph(map[string]*Node{"a": &nodeAWithSingleEdgeB}),
			node:    &nodeAWithSingleEdgeB,
			want:    Graph(map[string]*Node{"a": &nodeAWithSingleEdgeB}),
			wantErr: nil,
		},

		{
			name:    "add-valid-node-to-non-empty-graph-should-return-graph-with-new-node-added",
			graph:   Graph(map[string]*Node{"a": &nodeAWithSingleEdgeB}),
			node:    &nodeAWithSingleEdgeB,
			want:    Graph(map[string]*Node{"c": &nodeCWithSingleEdgeB}),
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := test.graph.AddNode(test.node)

			errorCmp := cmp.Comparer(func(x, y error) bool {
				if x != nil && y == nil {
					return false
				}

				if x == nil && y != nil {
					return false
				}

				if x.Error() != y.Error() {
					return false
				}
				return true
			})

			if ok := cmp.Equal(gotErr, test.wantErr, errorCmp); !ok {
				t.Errorf("(%+v) AddNode(%+v) = %v wantErr %v", test.graph, test.node, gotErr, test.wantErr)
			}

			if diff := cmp.Diff(test.want, test.graph); diff != "" {
				t.Errorf("(%+v) AddNode(%+v) mismatch (-want +got):\n%s", test.graph, test.node, diff)
			}
		})
	}
}

func TestContainsNode(t *testing.T) {
	tests := []struct {
		name  string
		graph Graph
		id    string
		want  bool
	}{
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.graph.ContainsNode(test.id)

			if got != test.want {
				t.Errorf("(%+v) ContainsNode(%v) = %v mismatch \n%v", test.graph, test.id, got, test.want)
			}
		})
	}
}
