package graph

import (
	"reflect"
	"testing"
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
			node:   &Node{"a", nil, map[string]WeightedEdge{}},
			dest:   &Node{"b", nil, map[string]WeightedEdge{}},
			weight: 1.0,
			want: Node{"a", nil, map[string]WeightedEdge{"b": WeightedEdge{
				Weight: 1.0,
				Dest:   &Node{ID: "b", Object: nil, Edges: map[string]WeightedEdge{}},
			}}},
		},

		{
			name:   "add-edge-to-same-node-does-nothing",
			node:   &Node{"a", nil, map[string]WeightedEdge{}},
			dest:   &Node{"a", nil, map[string]WeightedEdge{}},
			weight: 1.0,
			want:   Node{"a", nil, map[string]WeightedEdge{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.AddEdge(tt.dest, tt.weight)

			if ok := reflect.DeepEqual(*tt.node, tt.want); !ok {
				t.Errorf("(%+v) AddEdge(%+v, %f) - mismatch \n%+v", tt.node, tt.dest, tt.weight, tt.want)
			}
		})
	}
}
