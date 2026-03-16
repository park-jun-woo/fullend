//ff:func feature=orchestrator type=formatter control=iteration
//ff:what PrintStatus writes the status lines to a writer.

package orchestrator

import (
	"fmt"
	"io"
)

// PrintStatus writes the status lines to w.
func PrintStatus(w io.Writer, lines []StatusLine) {
	if len(lines) == 0 {
		fmt.Fprintln(w, "No SSOTs found.")
		return
	}

	fmt.Fprintln(w, "SSOT Status:")
	for _, l := range lines {
		fmt.Fprintf(w, "  %-12s %-30s %s\n", l.Kind, l.Path, l.Summary)
	}
}
