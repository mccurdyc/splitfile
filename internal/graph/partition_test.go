package graph

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_bfs(t *testing.T) {
	var tests = map[string]struct {
		input map[string]*Node
		want  map[string]*Node
	}{
		"node_two_changed_children_should_return_two_next-level_nodes": {
			input: map[string]*Node{"a": &Node{
				ID: "a",
				Edges: map[string]WeightedEdge{
					"b": {
						Weight: 1.0,
						Source: &Node{ID: "a", ShortestPath: 1.0},
						Dest:   &Node{ID: "b", ShortestPath: 3.0, ShortestPaths: []float64{3.0}, Edges: map[string]WeightedEdge{}},
					},
					"c": {
						Weight: 4.0,
						Source: &Node{ID: "a", ShortestPath: 1.0},
						Dest:   &Node{ID: "c", ShortestPath: 6.0, ShortestPaths: []float64{6.0}, Edges: map[string]WeightedEdge{}},
					},
				},
				ShortestPath:  1.0,
				ShortestPaths: []float64{1.0},
			}},
			want: map[string]*Node{
				"b": &Node{
					ID:            "b",
					Edges:         map[string]WeightedEdge{},
					ShortestPath:  2.0,
					ShortestPaths: []float64{3.0, 2.0},
				},
				"c": &Node{
					ID:            "c",
					Edges:         map[string]WeightedEdge{},
					ShortestPath:  5.0,
					ShortestPaths: []float64{6.0, 5.0},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := bfs(tt.input)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("bfs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_shortestPath(t *testing.T) {
	type input struct {
		edge WeightedEdge
	}

	type want struct {
		shortest float64
		changed  bool
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"nil_source_node_with_edge_weight_less_than_current": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: nil,
					Dest:   &Node{ShortestPath: 4.0},
				},
			},
			want: want{
				shortest: 2.0,
				changed:  true,
			},
		},

		"nil_source_node_with_edge_weight_greater_than_current": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: nil,
					Dest:   &Node{ShortestPath: 1.0},
				},
			},
			want: want{
				shortest: 1.0,
				changed:  false,
			},
		},

		"source_node_with_default_weight": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: &Node{ShortestPath: defaultWeight},
				},
			},
			want: want{
				shortest: 2.0,
				changed:  true,
			},
		},

		"nil_dest_node_new_shortest": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: &Node{ShortestPath: 1.0},
					Dest:   nil,
				},
			},
			want: want{
				shortest: 3.0,
				changed:  true,
			},
		},

		"dest_node_with_default_weight_force_new_shortest": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: &Node{ShortestPath: 3.0},
					Dest:   &Node{ShortestPath: defaultWeight},
				},
			},
			want: want{
				shortest: 5.0,
				changed:  true,
			},
		},

		"new_shortest_path": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: &Node{ShortestPath: 2.0},
					Dest:   &Node{ShortestPath: 5.0},
				},
			},
			want: want{
				shortest: 4.0,
				changed:  true,
			},
		},

		"equal_shortest_path": {
			input: input{
				edge: WeightedEdge{
					Weight: 2.0,
					Source: &Node{ShortestPath: 2.0},
					Dest:   &Node{ShortestPath: 4.0},
				},
			},
			want: want{
				shortest: 4.0,
				changed:  false,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotShortest, gotChanged := shortestPath(tt.input.edge)

			if gotShortest != tt.want.shortest {
				t.Errorf("shortestPath() = %+v - want: %+v\n\tinput: %+v\n", gotShortest, tt.want.shortest, tt.input)
			}

			if gotChanged != tt.want.changed {
				t.Errorf("shortestPath() = %+v - want: %+v\n\tinput: %+v\n", gotChanged, tt.want.changed, tt.input)
			}
		})
	}
}
