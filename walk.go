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

type FileInfo struct {
	path string
	err  error
}

func files(root string) <-chan FileInfo {
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
		err := filepath.Walk(root, walker)
		if err != nil {
			log.Panic(err) // Walker never returns error.
		}
	}()

	return ch
}

type LatLong struct {
	path      string
	latitude  float64
	longitude float64
	err       error
}

func extract(file string) (ll *LatLong) {
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

	return &LatLong{
		path:      file,
		latitude:  gi.Latitude.Decimal(),
		longitude: gi.Longitude.Decimal(),
	}
}

func latlong(files <-chan FileInfo) <-chan *LatLong {
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

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func csv(latlongs <-chan *LatLong, outf io.Writer, errf io.Writer) error {
	var err error
	for ll := range latlongs {
		if ll.err != nil {
			fmt.Fprintf(errf, "%s\n", ll.err)
			if err == nil {
				err = ll.err
			}
		} else {
			fmt.Fprintf(outf, "%q,%v,%v\n", ll.path, ll.latitude, ll.longitude)
		}
	}
	return err
}

func walk(root string, outf io.Writer, errf io.Writer) error {
	err := csv(latlong(files(root)), outf, errf)

	// Failing on `root` is pretty bad, anything else we can handle and consider
	// success.
	if e, ok := err.(*os.PathError); ok && e.Path == root {
		return err
	}

	return nil
}
