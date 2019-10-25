package splitfile

import (
	"go/token"
	"go/types"
	"testing"
)

type mockObject struct {
	pkgFn  func() *types.Package
	nameFn func() string
	posFn  func() token.Pos
}

func (mo mockObject) Pkg() *types.Package { return mo.pkgFn() }
func (mo mockObject) Name() string        { return mo.nameFn() }
func (mo mockObject) Pos() token.Pos      { return mo.posFn() }

func TestId(t *testing.T) {
	var tests = []struct {
		name   string
		want   string
		pkgFn  func() *types.Package
		nameFn func() string
		posFn  func() token.Pos
	}{
		{
			name:   "no package",
			want:   "name 1",
			pkgFn:  func() *types.Package { return nil },
			nameFn: func() string { return "name" },
			posFn:  func() token.Pos { return token.Pos(1) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := mockObject{
				pkgFn:  tt.pkgFn,
				nameFn: tt.nameFn,
				posFn:  tt.posFn,
			}

			got := Id(o)
			if got != tt.want {
				t.Errorf("Id(): want '%s', got '%s'", tt.want, got)
			}
		})
	}
}
