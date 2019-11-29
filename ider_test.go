package splitfile

import (
	"go/types"
	"testing"
)

type mockType struct {
	s string
}

func (mt mockType) Underlying() types.Type { return nil }
func (mt mockType) String() string         { return mt.s }

type mockObject struct {
	pkgFn  func() *types.Package
	nameFn func() string
	typeFn func() types.Type
}

func (mo mockObject) Pkg() *types.Package { return mo.pkgFn() }
func (mo mockObject) Name() string        { return mo.nameFn() }
func (mo mockObject) Type() types.Type    { return mo.typeFn() }

func TestId(t *testing.T) {
	var tests = []struct {
		name   string
		want   string
		pkgFn  func() *types.Package
		nameFn func() string
		typeFn func() types.Type
	}{
		{
			name:   "no_package",
			want:   "name type",
			pkgFn:  func() *types.Package { return nil },
			nameFn: func() string { return "name" },
			typeFn: func() types.Type { return mockType{s: "type"} },
		},

		{
			name: "package_exists",
			want: "package pkg (\"github.com/mccurdyc/path\") name type",
			pkgFn: func() *types.Package {
				return types.NewPackage("github.com/mccurdyc/path", "pkg")
			},
			nameFn: func() string { return "name" },
			typeFn: func() types.Type { return mockType{s: "type"} },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := mockObject{
				pkgFn:  tt.pkgFn,
				nameFn: tt.nameFn,
				typeFn: tt.typeFn,
			}

			got := Id(o)
			if got != tt.want {
				t.Errorf("Id(): want '%s', got '%s'", tt.want, got)
			}
		})
	}
}
