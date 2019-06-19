package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWalkNexist(t *testing.T) {
	var b strings.Builder
	var e strings.Builder
	err := walk("no-such-dir", &b, &e)
	require.Error(t, err)
	t.Logf("out: %q", e.String())
	require.Empty(t, b.String())
	require.Contains(t, e.String(), "no-such-dir: no such file or directory")
}

func TestWalk(t *testing.T) {
	var b strings.Builder
	var e strings.Builder
	err := walk("testdata", &b, &e)
	require.NoError(t, err)
	t.Logf("out: %q", b.String())
	require.Contains(t, b.String(),
		`"testdata/anubis.jpg",49.254444444444445,-123.1`)
	require.Contains(t, b.String(),
		`"testdata/subdir/wax-card.jpg",49.254444444444445,-123.1`)
	require.Empty(t, e.String())
}
