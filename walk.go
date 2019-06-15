package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func walk(root string, out io.Writer) error {
	fmt.Fprintf(out, "tree-of-images: %s\n", root)

	return filepath.Walk(root, walker)
}

func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return nil
	}

	fmt.Println(info.Name())
	return nil
}
