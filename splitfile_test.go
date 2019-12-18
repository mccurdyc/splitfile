package splitfile

import (
	"github.com/mccurdyc/splitfile/internal/graph"
)

var (
	emptyNode = graph.Node{
		ID:            "node",
		Object:        nil,
		Edges:         make(map[string]graph.WeightedEdge),
		Parents:       make(map[string]graph.WeightedEdge),
		ShortestPath:  -1.0,
		ShortestPaths: make([]float64, 0),
	}
)

// func Test_addRelated(t *testing.T) {
// 	var tests = []struct {
// 		name            string
// 		graph           graph.Graph
// 		pkgpath         string
// 		files           map[string]string
// 		node            *graph.Node
// 		typeNameToCheck string
// 		want            graph.Graph
// 		wantErr         error
// 	}{
// 		{
// 			name:    "type with no methods",
// 			graph:   make(graph.Graph),
// 			pkgpath: "a",
// 			files: map[string]string{"a/a.go": `
// 			package a
// 			type a int
// 			`},
// 			node:            &emptyNode,
// 			typeNameToCheck: "a",
// 			want:            make(graph.Graph),
// 			wantErr:         nil,
// 		},
//
// 		{
// 			name:    "type with a single method",
// 			graph:   make(graph.Graph),
// 			pkgpath: "a",
// 			files: map[string]string{"a/a.go": `
// 			package a
// 			type a int
// 			func (a a) Val() int {
// 				return int(a)
// 			}
// 			`},
// 			node: &graph.Node{
// 				ID: "package a (\"a\")",
// 			},
// 			typeNameToCheck: "a",
// 			want: graph.Graph{
// 				"package a (\"a\")": &graph.Node{},
// 				"package a (\"a\") Val func() int": &graph.Node{
// 					ID:            "package a (\"a\") Val func() int",
// 					Object:        nil,
// 					Edges:         map[string]graph.WeightedEdge{},
// 					Parents:       map[string]graph.WeightedEdge{},
// 					ShortestPath:  -1.0,
// 					ShortestPaths: []float64{},
// 				},
// 				"package a (\"a\") a a.a": &graph.Node{
// 					ID:            "package a (\"a\") a a.a",
// 					Object:        nil,
// 					Edges:         map[string]graph.WeightedEdge{},
// 					Parents:       map[string]graph.WeightedEdge{},
// 					ShortestPath:  -1.0,
// 					ShortestPaths: []float64{},
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			fset := token.NewFileSet()
// 			var files []*ast.File
//
// 			for fname, fileContents := range tt.files {
// 				f, err := parser.ParseFile(fset, fname, fileContents, 0)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
//
// 				files = append(files, f)
// 			}
//
// 			conf := types.Config{Importer: importer.Default()}
// 			pkg, err := conf.Check(tt.pkgpath, fset, files, nil)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			tt.node.Object = pkg.Scope().Lookup(tt.typeNameToCheck)
//
// 			gotErr := addRelated(tt.graph, tt.node)
//
// 			if !reflect.DeepEqual(tt.want, tt.graph) {
// 				t.Errorf("addRelated() mismatch:\n\twant: %+v\n\tgot: %+v", tt.want, tt.graph)
// 			}
//
// 			// https://github.com/google/go-cmp/issues/24
// 			errorCmp := func(x, y error) bool {
// 				if x == nil || y == nil {
// 					return x == nil && y == nil
// 				}
// 				return x.Error() == y.Error()
// 			}
//
// 			if ok := errorCmp(gotErr, tt.wantErr); !ok {
// 				t.Errorf("addRelated() = %v, wantErr %v", gotErr, tt.wantErr)
// 			}
// 		})
// 	}
// }
