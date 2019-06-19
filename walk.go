package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	exif "github.com/dsoprea/go-exif"
)

func files(root string) (<-chan string, <-chan error) {
	ch := make(chan string)
	cherr := make(chan error)

	var walker = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			ch <- path
		}

		return nil
	}

	go func() {
		defer close(ch)
		defer close(cherr)
		err := filepath.Walk(root, walker)
		if err != nil {
			cherr <- err
		}
	}()

	return ch, cherr
}

type LatLong struct {
	path      string
	latitude  float64
	longitude float64
}

func extract(file string) (ll LatLong, ok bool) {
	// Despite the declaration, exif.SearchFileAndExtractExif never returns
	// a non-nil error, it panics if anything goes wrong! It shouldn't do
	// this, its unidiomatic AFAICT.  I could open the file myself, slurp it
	// in, and call SearchAndExtractExif (which doesn't panic), but recover
	// also works, and avoids having to copy-n-paste in the file slurping
	// code.
	defer func() {
		recover()
	}()
	rawExif, err := exif.SearchFileAndExtractExif(file)
	if err != nil {
		return
	}
	im := exif.NewIfdMapping()

	err = exif.LoadStandardIfds(im)
	if err != nil {
		log.Panic(err) // go-exif can't load its internals, just die.
	}

	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return // Stuff that looked like exif wasn't, ignore.
	}

	ifd, err := index.RootIfd.ChildWithIfdPath(exif.IfdPathStandardGps)
	if err != nil {
		return // No GPS info, ignore
	}

	gi, err := ifd.GpsInfo()
	if err != nil {
		return // Stuff that is supposed to be GPS info wasn't, ignore.
	}

	return LatLong{file, gi.Latitude.Decimal(), gi.Longitude.Decimal()}, true
}

func latlong(files <-chan string, errors <-chan error) (
	<-chan LatLong, <-chan error) {
	ch := make(chan LatLong)
	cherr := make(chan error)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for {
			select {
			case err := <-errors:
				cherr <- err
				return // Fail fast and pass it on
			case file, ok := <-files:
				if !ok {
					return
				}
				if ll, ok := extract(file); ok {
					ch <- ll
				}
			}
		}
	}

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		wg.Wait()
		close(ch)
		close(cherr)
	}()

	return ch, cherr
}

func csv(ll <-chan LatLong, errors <-chan error, out io.Writer) error {
	for {
		select {
		case err := <-errors:
			if err != nil {
				return err
			}
		case exif, ok := <-ll:
			if !ok {
				return nil
			}
			fmt.Fprintf(out, "%q,%v,%v\n",
				exif.path, exif.latitude, exif.longitude)
		}
	}
}

func walk(root string, out io.Writer) error {
	f, cherr := files(root)
	l, cherr := latlong(f, cherr)
	return csv(l, cherr, out)
}
