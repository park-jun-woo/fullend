//ff:func feature=reporter type=formatter control=iteration dimension=1
//ff:what 검증 결과 요약 메시지를 출력한다
package reporter

import (
	"fmt"
	"io"
)

// printSummary writes the final summary line of the report to w.
func printSummary(w io.Writer, r *Report) {
	fmt.Fprintln(w)

	if r.HasFailure() {
		fmt.Fprintln(w, "FAILED: Fix errors before codegen.")
		return
	}

	allSkip := true
	for _, s := range r.Steps {
		if s.Status == Pass {
			allSkip = false
			break
		}
	}
	if allSkip {
		fmt.Fprintln(w, "No SSOT sources found.")
		return
	}

	hasSkip := false
	for _, s := range r.Steps {
		if s.Status == Skip {
			hasSkip = true
			break
		}
	}
	if hasSkip {
		fmt.Fprintln(w, "Partial validation passed.")
	} else {
		fmt.Fprintln(w, "All SSOT sources are consistent.")
	}
}
