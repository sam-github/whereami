package main

import (
	"fmt"
	"html"
	"io"
)

var BEG = `<!DOCTYPE html>
<html>
<head><title>LatLong from files in %q</title></head>
<body>
<h1>LatLong from files in %q</h1>
<ul>
`

var END = `
</ul>
</body>
</html>
`

func Html(root string, outf io.Writer) Lister {
	fmt.Fprintf(outf, BEG, root, root)

	return func(ch <-chan *LatLongInfo, outf io.Writer, errf io.Writer) error {
		var err error

		for ll := range ch {
			if ll.err != nil {
				fmt.Fprintf(errf, "%s\n", ll.err)
				if err == nil {
					err = ll.err
				}
			} else {
				fmt.Fprintf(outf, "<li><tt>%s</tt>: %v,%v</li>\n",
					html.EscapeString(ll.path), ll.latitude, ll.longitude)
			}
		}
		fmt.Fprintf(outf, "%s", END)

		return err
	}
}
