//ff:func feature=cli type=formatter control=selection
//ff:what history 결과를 yaml/json 포맷으로 출력
package main

import (
	"fmt"
	"io"
	"time"

	"github.com/clari/whyso/pkg/history"
)

func formatHistory(w io.Writer, h *history.FileHistory, format string) {
	switch format {
	case "json":
		fmt.Fprintf(w, "{\n")
		fmt.Fprintf(w, "  \"apiVersion\": \"whyso/v1\",\n")
		fmt.Fprintf(w, "  \"file\": %q,\n", h.File)
		fmt.Fprintf(w, "  \"created\": %q,\n", h.Created.Format(time.RFC3339))
		fmt.Fprintf(w, "  \"history\": []\n") // simplified
		fmt.Fprintf(w, "}\n")
	default:
		fmt.Fprintf(w, "apiVersion: whyso/v1\n")
		fmt.Fprintf(w, "file: %s\n", h.File)
		fmt.Fprintf(w, "created: %s\n", h.Created.Format(time.RFC3339))
		fmt.Fprintf(w, "history:\n")
		formatHistoryEntries(w, h.History)
	}
}
