//ff:func feature=reporter type=formatter control=selection
//ff:what 검증 단계 하나를 상태에 따라 포맷팅하여 출력한다
package reporter

import (
	"fmt"
	"io"
)

// printStep writes a single step result to w.
func printStep(w io.Writer, step StepResult) {
	// Separator step.
	if step.Name == "---" {
		fmt.Fprintf(w, "\n── %s ──\n", step.Summary)
		return
	}

	switch step.Status {
	case Pass:
		fmt.Fprintf(w, "✓ %-12s %s\n", step.Name, step.Summary)
		printErrors(w, step.Errors, step.Suggestions)
	case Fail:
		fmt.Fprintf(w, "✗ %-12s %s\n", step.Name, step.Summary)
		printErrors(w, step.Errors, step.Suggestions)
	case Skip:
		fmt.Fprintf(w, "— %-12s %s\n", step.Name, step.Summary)
	}
}
