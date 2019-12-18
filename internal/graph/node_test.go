package graph

import (
	"reflect"
	"testing"
)

var (
	nodeA = Node{
		ID:            "a",
		Object:        nil,
		Edges:         map[string]WeightedEdge{},
		Parents:       map[string]WeightedEdge{},
		ShortestPath:  1.0,
		ShortestPaths: []float64{1.0},
	}

	nodeB = Node{
		ID:            "b",
		Object:        nil,
		Edges:         map[string]WeightedEdge{},
		Parents:       map[string]WeightedEdge{},
		ShortestPath:  1.0,
		ShortestPaths: []float64{1.0},
	}
)

func TestAddEdge(t *testing.T) {
	tests := []struct {
		name   string
		node   *Node
		dest   *Node
		weight float64
		want   Node
	}{
		{
			name:   "add-edge-w1.0-to-empty-slice-should-return-edges-with-single-value",
			node:   &nodeA,
			dest:   &nodeB,
			weight: 1.0,
			want: Node{
				ID:     "a",
				Object: nil,
				Edges: map[string]WeightedEdge{
					"b": WeightedEdge{
						Weight: 1.0,
						Source: &nodeA,
						Dest:   &nodeB,
					},
				},
				Parents:       map[string]WeightedEdge{},
				ShortestPath:  1.0,
				ShortestPaths: []float64{1.0},
			},
		},

		{
			name:   "add-edge-to-same-node-does-nothing",
			node:   &nodeA,
			dest:   &nodeA,
			weight: 1.0,
			want:   nodeA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.AddEdge(tt.dest, tt.weight)

			if ok := reflect.DeepEqual(*tt.node, tt.want); !ok {
				t.Errorf("AddEdge() - mismatch \n\twant: %+v\n\tgot:%+v", tt.want, *tt.node)
			}
		})
	}
}
