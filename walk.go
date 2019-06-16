package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	exif "github.com/dsoprea/go-exif"
)

func walk(root string, out io.Writer) error {
	var walker = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		// Despite the declaration, this never return a non-nil error, it panics
		// if anything goes wrong! It shouldn't do this, its unidiomatic AFAICT.
		// I could open the file myself, slurp it in, and call
		// SearchAndExtractExif (which doesn't panic), but recover also works,
		// and avoids having to copy-n-paste in the file slurping code.
		defer func() {
			recover()
		}()
		rawExif, err := exif.SearchFileAndExtractExif(path)
		if err != nil {
			return nil // File had no exif, ignore.
		}
		im := exif.NewIfdMapping()

		err = exif.LoadStandardIfds(im)
		if err != nil {
			log.Panic(err) // go-exif can't load its internals, fatal.
		}

		ti := exif.NewTagIndex()

		_, index, err := exif.Collect(im, ti, rawExif)
		if err != nil {
			return nil // Stuff that looked like exif wasn't, ignore.
		}

		ifd, err := index.RootIfd.ChildWithIfdPath(exif.IfdPathStandardGps)
		if err != nil {
			return nil // No GPS info, ignore
		}

		gi, err := ifd.GpsInfo()
		if err != nil {
			return nil // Stuff that is supposed to be GPS info wasn't, ignore.
		}

		fmt.Fprintf(out, "%q,%v,%v\n",
			path, gi.Latitude.Decimal(), gi.Longitude.Decimal())

		return nil
	}
	return filepath.Walk(root, walker)
}
