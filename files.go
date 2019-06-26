package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/MichaelTJones/walk"
)

// Unfortunately, instead of using filepath.WalkFunc as its callback type, walk
// uses a custom typedef, making the type signature of walk.Walk() gratuitously
// incompatible with filepath.Walk(), so make a wrapper.
func _walk(root string, cb filepath.WalkFunc) error {
	var wrap = func(path string, info os.FileInfo, err error) error {
		return cb(path, info, err)
	}
	return walk.Walk(root, wrap)
}

type FileInfo struct {
	path string
	err  error
}

// Use stdlib's pathwalk() with no parallelism.
func Files(root string) <-chan FileInfo {
	return _files(filepath.Walk, root)
}

func FilesMtj(root string) <-chan FileInfo {
	return _files(_walk, root)
}

type walkFn func(string, filepath.WalkFunc) error

func _files(wfn walkFn, root string) <-chan FileInfo {
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
		err := wfn(root, walker)
		if err != nil {
			log.Panic(err) // walker() never returns error.
		}
	}()

	return ch
}

func FilesJ(j int, root string) <-chan FileInfo {
	if j < 1 {
		j = runtime.GOMAXPROCS(0)
	}

	ch := make(chan FileInfo)

	dirs := make(chan string)
	var wg sync.WaitGroup

	var walker = func() {
		for path := range dirs {
			readdir(&wg, path, dirs, ch)
		}
	}

	for i := 0; i < j; i++ {
		go walker()
	}

	go func() {
		defer close(ch)
		wg.Add(1)
		dirs <- root
		wg.Wait()
		close(dirs)
	}()

	return ch

}

func readdir(
	wg *sync.WaitGroup, path string, dirs chan string, ch chan<- FileInfo,
) {
	defer wg.Done()

	files, err := ioutil.ReadDir(path)

	if err != nil {
		ch <- FileInfo{path, err}
		return
	}

	for _, fi := range files {
		p := path + "/" + fi.Name()
		if fi.Mode().IsRegular() {
			ch <- FileInfo{p, nil}
		} else {
			wg.Add(1)
			// Writing will block until it is read, but since we are the reader,
			// this could deadlock. The workers should probably be working on
			// some kind of list, not a channel.
			go func() {
				dirs <- p
			}()
		}
	}
}
