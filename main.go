package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var HELP = `usage: %s [-j N] [-h] <tree-of-images>
  Print to stdout location information for all images found
  in 'tree-of-images'.

`

func main() {
	j := flag.Int("j", runtime.GOMAXPROCS(0), "use `N` concurrent workers")
	h := flag.Bool("h", false, "print a helpful message")

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

	err := walk(*j, flag.Arg(0), os.Stdout, os.Stderr)

	if err != nil {
		os.Exit(1)
	}
}
