package main

import (
	"fmt"
	"io"
)

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
