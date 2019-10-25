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
	if ider.Pkg() != nil {
		return fmt.Sprintf("%s %s %d", ider.Pkg().String(), ider.Name(), ider.Pos())
	}

	return fmt.Sprintf("%s %d", ider.Name(), ider.Pos())
}
