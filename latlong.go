package main

import (
	"runtime"
	"sync"
)

func latlong(j int, files <-chan FileInfo) <-chan *LatLong {
	if j < 1 {
		j = runtime.GOMAXPROCS(0)
	}
	ch := make(chan *LatLong)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for fi := range files {
			if fi.err != nil {
				ch <- &LatLong{path: fi.path, err: fi.err}
			} else if ll := extract(fi.path); ll != nil {
				ch <- ll
			}
		}
	}

	for i := 0; i < j; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
