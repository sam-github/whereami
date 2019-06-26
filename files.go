package main

import (
	"log"
	"os"

	"github.com/MichaelTJones/walk"
)

type FileInfo struct {
	path string
	err  error
}

func Files(root string) <-chan FileInfo {
	ch := make(chan FileInfo)

	var walker = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			ch <- FileInfo{path, err}
		} else if info.Mode().IsRegular() {
			ch <- FileInfo{path, nil}
		}

		return nil
	}

	go func() {
		defer close(ch)
		err := walk.Walk(root, walker)
		if err != nil {
			log.Panic(err) // walker() never returns error.
		}
	}()

	return ch
}
