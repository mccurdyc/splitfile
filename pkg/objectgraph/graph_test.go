package objectgraph

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name              string
		expectedRelations map[types.Object][]types.Object
	}{
		{
			name:              "empty relations map",
			expectedRelations: map[types.Object][]types.Object{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := New()

			assert.Equal(t, test.expectedRelations, actual.relations, "relations map does not match expected.")
		})
	}
}
