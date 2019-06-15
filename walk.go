package main

import (
	"fmt"
	"io"
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

		rawExif, err := exif.SearchFileAndExtractExif(path)
		if err != nil {
			return err
		}
		im := exif.NewIfdMapping()

		err = exif.LoadStandardIfds(im)
		if err != nil {
			return err
		}

		ti := exif.NewTagIndex()

		_, index, err := exif.Collect(im, ti, rawExif)
		if err != nil {
			return err
		}

		ifd, err := index.RootIfd.ChildWithIfdPath(exif.IfdPathStandardGps)
		if err != nil {
			return err
		}

		gi, err := ifd.GpsInfo()
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "%q,%v,%v\n",
			path, gi.Latitude.Decimal(), gi.Longitude.Decimal())

		return nil
	}
	return filepath.Walk(root, walker)
}
