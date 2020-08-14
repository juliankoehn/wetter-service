package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAtomicBool(t *testing.T) {
	ab := NewAtomicBool(false)
	require.False(t, ab.Get())

	ab = NewAtomicBool(true)
	require.True(t, ab.Get())

	require.True(t, ab.Set(false))
	require.False(t, ab.Get())

	require.False(t, ab.Set(false))
	require.False(t, ab.Set(false))
	require.False(t, ab.Get())

	require.False(t, ab.Set(true))
	require.True(t, ab.Get())
}
