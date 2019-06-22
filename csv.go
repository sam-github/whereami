package main

import (
	"fmt"
	"io"
)

func Csv() Lister {
	return func(
		ch <-chan *LatLongInfo, outf io.Writer, errf io.Writer,
	) error {
		var err error
		for ll := range ch {
			if ll.err != nil {
				fmt.Fprintf(errf, "%s\n", ll.err)
				if err == nil {
					err = ll.err
				}
			} else {
				fmt.Fprintf(outf, "%q,%v,%v\n",
					ll.path, ll.latitude, ll.longitude)
			}
		}
		return err
	}
}
