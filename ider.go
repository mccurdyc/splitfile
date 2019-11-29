package splitfile

import (
	"fmt"
	"go/types"
)

type Ider interface {
	Pkg() *types.Package
	Name() string
	Type() types.Type
}

func Id(ider Ider) string {
	id := fmt.Sprintf("%s %s", ider.Name(), ider.Type().String())

	if ider.Pkg() != nil {
		id = fmt.Sprintf("%s %s", ider.Pkg().String(), id)
	}

	return id
}
