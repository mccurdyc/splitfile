package splitfile

import (
	"fmt"
	"go/token"
	"go/types"
)

type Ider interface {
	Pkg() *types.Package
	Name() string
	Pos() token.Pos
}

func Id(ider Ider) string {
	var pkg string
	if ider.Pkg() != nil {
		pkg = ider.Pkg().String()
	}

	return fmt.Sprintf("%q %q %d", pkg, ider.Name(), ider.Pos())
}
