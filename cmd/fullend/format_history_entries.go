//ff:func feature=cli type=formatter control=iteration dimension=1
//ff:what history 엔트리 목록을 YAML 형식으로 출력
package main

import (
	"fmt"
	"io"
	"time"

	"github.com/clari/whyso/pkg/history"
)

func formatHistoryEntries(w io.Writer, entries []history.ChangeEntry) {
	for _, e := range entries {
		fmt.Fprintf(w, "  - timestamp: %s\n", e.Timestamp.Format(time.RFC3339))
		fmt.Fprintf(w, "    session: %s\n", e.Session)
		fmt.Fprintf(w, "    user_request: %q\n", e.UserRequest)
		if e.Answer != "" {
			fmt.Fprintf(w, "    answer: %q\n", e.Answer)
		}
		fmt.Fprintf(w, "    tool: %s\n", e.Tool)
		if e.Subagent {
			fmt.Fprintf(w, "    subagent: true\n")
		}
		formatHistoryEntrySources(w, e.Sources)
	}
}
