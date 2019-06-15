package main

import (
	"fmt"
	"log"
	"os"
)

var HELP = `usage: %s <tree-of-images>
  Print to stdout location information for all images found
  in 'tree-of-images'.
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, HELP, os.Args[0])
		os.Exit(2)
	}
	if os.Args[1] == "-h" || os.Args[1] == "-help" {
		fmt.Fprintf(os.Stdout, HELP, os.Args[0])
		return
	}

	err := walk(os.Args[1], os.Stdout)

	if err == nil {
		return
	}

	log.Fatal(err)
}
