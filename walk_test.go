package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWalkNexist(t *testing.T) {
	var b strings.Builder
	var e strings.Builder
	err := walk(0, "no-such-dir", &b, &e)
	require.Error(t, err)
	t.Logf("out: %q", e.String())
	require.Empty(t, b.String())
	require.Contains(t, e.String(), "no-such-dir: no such file or directory")
}

func TestWalk(t *testing.T) {
	var b strings.Builder
	var e strings.Builder
	err := walk(0, "testdata", &b, &e)
	require.NoError(t, err)
	t.Logf("out: %q", b.String())
	require.Contains(t, b.String(),
		`"testdata/anubis.jpg",49.254444444444445,-123.1`)
	require.Contains(t, b.String(),
		`"testdata/subdir/wax-card.jpg",49.254444444444445,-123.1`)
	require.Empty(t, e.String())
}

// There isn't enough data here to benchmark, so set BENCHROOT in the
// environment to your image set.
var BENCHROOT string = "testdata"

func init() {
	if br := os.Getenv("BENCHROOT"); len(br) > 0 {
		BENCHROOT = br
	}
}

func BenchmarkFilepath_Parallel_Csv(b *testing.B) {
	d := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := csv(latlong(0, files(BENCHROOT)), d, d)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkFilepath_Sequential_Csv(b *testing.B) {
	d := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := csv(latlong(1, files(BENCHROOT)), d, d)
		if err != nil {
			b.FailNow()
		}
	}
}
