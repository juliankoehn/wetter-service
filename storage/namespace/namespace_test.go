package namespace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNameSpace(t *testing.T) {
	SetNamespace("gotest")

	ns := GetNamespace()
	require.Equal(t, "gotest", ns)
}
