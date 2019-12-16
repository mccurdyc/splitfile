package graph

import "testing"

func Test_shortestPath(t *testing.T) {
	type input struct {
		edge WeightedEdge
	}

	tests := map[string]struct {
		input input
		want  float64
	}{
		"nil_source_node_with_edge_weight_2.0": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: nil,
				},
			},
			want: 2.0,
		},

		"source_node_with_default_weight_edge_weight_2.0": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: &Node{
						ShortestPath: defaultWeight,
					},
				},
			},
			want: 2.0,
		},

		"source_node_with_shortest_path_of_2.0_edge_weight_2.0": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: &Node{
						ShortestPath: 2.0,
					},
				},
			},
			want: 4.0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := shortestPath(tt.input.edge)

			if got != tt.want {
				t.Errorf("shortestPath() = %+v - want: %+v\n\tinput: %+v\n", got, tt.want, tt.input)
			}
		})
	}
}
