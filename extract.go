package main

import (
	"log"

	exif "github.com/dsoprea/go-exif"
)

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
