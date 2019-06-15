package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func walk(root string, out io.Writer) error {
	var walker = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		fmt.Fprintf(out, "%s\n", info.Name())
		return nil
	}
	return filepath.Walk(root, walker)
}
