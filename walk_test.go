package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWalkNexist(t *testing.T) {
	var b strings.Builder
	err := walk("no-such-dir", &b)
	require.Error(t, err)
	require.Empty(t, b.String())
}

func TestWalk(t *testing.T) {
	var b strings.Builder
	err := walk("testdata", &b)
	require.NoError(t, err)
	require.Contains(t, b.String(), "anubis.jpg")
	require.Contains(t, b.String(), "wax-card.jpg")
}
