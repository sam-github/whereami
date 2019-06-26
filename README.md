# whereami?

Given a directory of some images, that contains a sub-directory that contains
some images, create a command line utility that reads the EXIF data from the
images and writes the image path, latitude and longitude to file as a CSV. Use
Go routines and channels to do the reads concurrently if possible. For extra
credit, provide an option to write to HTML as well.

## Benchmarks

Run with:

	BENCHROOT=/path/to/deep/tree/of/images make bench

Unsurprisingly, the exif parsing benefits a lot from a goroutine pool, but
using a parallel fs walker helps a bit.

```
% BENCHROOT=~/Dropbox/Home/Photos make bench
go test -v -bench=.
BENCHROOT= /home/sam/Dropbox/Home/Photos
GOMAXPROCS= 4
=== RUN   TestWalkNexist
--- PASS: TestWalkNexist (0.00s)
    walk_test.go:19: out: "open no-such-dir: no such file or directory\n"
=== RUN   TestWalk
--- PASS: TestWalk (0.00s)
    walk_test.go:29: out: "\"testdata/subdir/wax-card.jpg\",49.254444444444445,-123.1\n\"testdata/anubis.jpg\",49.254444444444445,-123.1\n"                                  
goos: linux
goarch: amd64
pkg: github.com/sam-github/whereami
BenchmarkMTJ_MAX__EXIF_MAX-4                   1        32688491663 ns/op
BenchmarkFILES_MAX__EXIF_MAX-4                 1        37722251814 ns/op
BenchmarkFILES_1__EXIF_MAX-4                   1        42303744933 ns/op
BenchmarkFILES_MAX__EXIF_1-4                   1        82082289850 ns/op
BenchmarkFILES_1__EXIF_1-4                     1        72609082373 ns/op
PASS
ok      github.com/sam-github/whereami  267.476s
```
