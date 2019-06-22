package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

var HELP = `usage: %s [-j N] [-h] <tree-of-images>
  Print to stdout location information for all images found
  in 'tree-of-images'.

`

// A function that can list LatLong in some format.
type Lister func(<-chan *LatLong, io.Writer, io.Writer) error

func main() {
	h := flag.Bool("h", false, "print a helpful message")
	j := flag.Int("j", runtime.GOMAXPROCS(0), "use `N` concurrent workers")
	H := flag.Bool("html", false, "report location in HTML, not CSV")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, HELP, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {
		flag.CommandLine.SetOutput(os.Stdout)
		fmt.Fprintf(os.Stdout, HELP, os.Args[0])
		flag.PrintDefaults()
		return
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Missing argument, try `%s -h`", os.Args[0])
		os.Exit(2)
	}

	var out Lister

	if *H {
		out = Html(flag.Arg(0), os.Stdout)
	} else {
		out = csv
	}

	err := walk(*j, flag.Arg(0), out, os.Stdout, os.Stderr)

	if err != nil {
		os.Exit(1)
	}
}
