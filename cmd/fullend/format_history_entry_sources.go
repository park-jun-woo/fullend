//ff:func feature=cli type=formatter control=selection
//ff:what history 엔트리의 source 필드 출력
package main

import (
	"fmt"
	"io"

	"github.com/clari/whyso/pkg/history"
)

func formatHistoryEntrySources(w io.Writer, sources []history.Source) {
	switch len(sources) {
	case 0:
		// no sources
	case 1:
		fmt.Fprintf(w, "    source: %s:%d\n", sources[0].File, sources[0].Line)
	default:
		fmt.Fprintf(w, "    sources:\n")
		for _, s := range sources {
			fmt.Fprintf(w, "      - %s:%d\n", s.File, s.Line)
		}
	}
}
