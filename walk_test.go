package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWalkNexist(t *testing.T) {
	var b strings.Builder
	var e strings.Builder
	err := Walk(1, "no-such-dir", Csv(), &b, &e)
	require.Error(t, err)
	t.Logf("out: %q", e.String())
	require.Empty(t, b.String())
	require.Contains(t, e.String(), "no-such-dir: no such file or directory")
}

func TestWalk(t *testing.T) {
	var b strings.Builder
	var e strings.Builder
	err := Walk(0, "testdata", Csv(), &b, &e)
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

	fmt.Println("BENCHROOT=", BENCHROOT)
	fmt.Println("GOMAXPROCS=", runtime.GOMAXPROCS(0))
}

func BenchmarkMTJ_MAX__EXIF_MAX(b *testing.B) {
	d := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := Csv()(LatLong(0, FilesMtj(BENCHROOT)), d, d)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkFILES_MAX__EXIF_MAX(b *testing.B) {
	d := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := Csv()(LatLong(0, FilesJ(0, BENCHROOT)), d, d)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkFILES_1__EXIF_MAX(b *testing.B) {
	d := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := Csv()(LatLong(0, Files(BENCHROOT)), d, d)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkFILES_MAX__EXIF_1(b *testing.B) {
	d := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := Csv()(LatLong(1, FilesJ(0, BENCHROOT)), d, d)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkFILES_1__EXIF_1(b *testing.B) {
	d := ioutil.Discard

	for i := 0; i < b.N; i++ {
		err := Csv()(LatLong(1, Files(BENCHROOT)), d, d)
		if err != nil {
			b.FailNow()
		}
	}
}
