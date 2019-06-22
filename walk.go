package main

import (
	"io"
	"os"
)

func walk(
	j int,
	root string,
	list Lister,
	outf io.Writer,
	errf io.Writer,
) error {
	err := list(latlong(j, files(root)), outf, errf)

	// Failing on `root` is pretty bad, anything else we can handle and consider
	// success.
	if e, ok := err.(*os.PathError); ok && e.Path == root {
		return err
	}

	return nil
}
