package splitfile

import (
	"go/token"
	"go/types"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/packages"
)

type mockType struct {
	s string
}

func (mt *mockType) Underlying() types.Type {
	return mt
}

func (mt *mockType) String() string {
	return mt.s
}

func TestCheckMethods(t *testing.T) {
	tests := []struct {
		name     string
		pkgpath  string
		files    map[string]string
		expected map[string][]string
	}{
		{
			name:    "no methods",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
		type a int
		`,
			},
			expected: make(map[string][]string),
		},
		{
			name:    "method no params one result",
			pkgpath: "a",
			files: map[string]string{"a/a.go": `package a
			type a int
			type b int

			func (a a) ma() b {
			return b(1)
			}
		`,
			},
			expected: map[string][]string{"a.a": []string{"a.a", "a.b", "method (a.a) ma() a.b"}, "a.b": []string{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir, cleanup, err := analysistest.WriteFiles(test.files)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			cfg := &packages.Config{
				Mode:  packages.LoadAllSyntax,
				Dir:   dir,
				Tests: true,
				Env:   append(os.Environ(), "GOPATH="+dir, "GO111MODULE=off", "GOPROXY=off"),
			}
			pkgs, err := packages.Load(cfg, test.pkgpath)
			if err != nil {
				t.Fatal(err)
			}

			for _, p := range pkgs {
				if len(p.Errors) > 0 {
					t.Fatal(p.Errors[0])
				}
				defs := p.TypesInfo.Defs
				for _, node := range defs {
					// nil for a package definition
					if node == nil {
						continue
					}

					actual := checkMethods(types.NewMethodSet(node.Type()))
					nodeKey := node.Type().String()
					assert.ElementsMatch(t, test.expected[nodeKey], actual, "unexpected results from checkMethods.")
				}
			}
		})
	}

}

func TestCheckSignature(t *testing.T) {
	tests := []struct {
		name     string
		recv     *types.Var
		params   []*types.Var
		results  []*types.Var
		variadic bool
		expected []string
	}{
		{
			name:     "nil receiver nil params nil results",
			recv:     nil,
			params:   nil,
			results:  nil,
			variadic: false,
			expected: []string{},
		},
		{
			name:   "nil receiver nil params single result",
			recv:   nil,
			params: nil,
			results: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
			},
			variadic: false,
			expected: []string{"a"},
		},
		{
			name:   "nil receiver nil params multiple results",
			recv:   nil,
			params: nil,
			results: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "b"}),
			},
			variadic: false,
			expected: []string{"a", "b"},
		},
		{
			name: "nil receiver single param nil results",
			recv: nil,
			params: []*types.Var{
				types.NewParam(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
			},
			results:  nil,
			variadic: false,
			expected: []string{"a"},
		},
		{
			name: "nil receiver (func) multiple params nil results",
			recv: nil,
			params: []*types.Var{
				types.NewParam(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
				types.NewParam(token.Pos(1), &types.Package{}, "abc", &mockType{s: "b"}),
			},
			results:  nil,
			variadic: false,
			expected: []string{"a", "b"},
		},
		{
			name: "non-nil receiver params results",
			recv: types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
			params: []*types.Var{
				types.NewParam(token.Pos(1), &types.Package{}, "abc", &mockType{s: "b"}),
				types.NewParam(token.Pos(1), &types.Package{}, "abc", &mockType{s: "c"}),
			},
			results: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "d"}),
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "e"}),
			},
			variadic: false,
			expected: []string{"a", "b", "c", "d", "e"},
		},
		{
			name: "non-nil receiver variadic params results",
			recv: types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
			params: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "b"}),
				types.NewVar(token.Pos(1), &types.Package{}, "abc", types.NewSlice(&mockType{s: "c"})),
			},
			results: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "d"}),
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "e"}),
			},
			variadic: true,
			expected: []string{"a", "b", "c", "d", "e"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := checkSignature(types.NewSignature(test.recv, types.NewTuple(test.params...), types.NewTuple(test.results...), test.variadic))
			assert.Equal(t, test.expected, actual, "check signature returned unexpected results.")
		})
	}
}

func TestCheckVar(t *testing.T) {
	tests := []struct {
		name     string
		v        *types.Var
		expected string
	}{
		{
			name:     "nil var",
			v:        nil,
			expected: "",
		},
		{
			name:     "nil casted var",
			v:        types.NewVar(token.Pos(0), &types.Package{}, "", types.Type(nil)),
			expected: "",
		},
		{
			name:     "nil type var",
			v:        types.NewVar(token.Pos(0), &types.Package{}, "", &mockType{}),
			expected: "",
		},
		{
			name:     "nil type var with package",
			v:        types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{}),
			expected: "",
		},
		{
			name:     "valid var",
			v:        types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
			expected: "a",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := checkVar(test.v)
			assert.Equal(t, test.expected, actual, "check var returned unexpected results.")
		})
	}
}

func TestCheckTuple(t *testing.T) {
	tests := []struct {
		name     string
		vars     []*types.Var
		expected []string
	}{
		{
			name:     "empty tuple",
			vars:     []*types.Var{},
			expected: []string{},
		},
		{
			name: "single empty var",
			vars: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: ""}),
			},
			expected: []string{},
		},
		{
			name: "multiple empty vars",
			vars: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: ""}),
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: ""}),
			},
			expected: []string{},
		},
		{
			name: "single valid var",
			vars: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
			},
			expected: []string{"a"},
		},
		{
			name: "multiple valid vars",
			vars: []*types.Var{
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "a"}),
				types.NewVar(token.Pos(1), &types.Package{}, "abc", &mockType{s: "b"}),
			},
			expected: []string{"a", "b"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := checkTuple(types.NewTuple(test.vars...))
			assert.Equal(t, test.expected, actual, "check tuple returned unexpected results.")
		})
	}
}
