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
	t.Logf("out: %q", b.String())
	require.Contains(t, b.String(),
		`"testdata/anubis.jpg",49.254444444444445,-123.1`)
	require.Contains(t, b.String(),
		`"testdata/subdir/wax-card.jpg",49.254444444444445,-123.1`)
}
