//ff:func feature=reporter type=formatter control=iteration dimension=1
//ff:what 계약 상태 보고서를 출력한다
package reporter

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/fullend/internal/contract"
)

// PrintContract writes the contract status report to w.
func PrintContract(w io.Writer, funcs []contract.FuncStatus) {
	if len(funcs) == 0 {
		fmt.Fprintln(w, "No contract directives found.")
		return
	}

	fmt.Fprintln(w, "Contract Status:")
	for _, f := range funcs {
		var icon string
		switch f.Status {
		case "gen":
			icon = " "
		case "preserve":
			icon = "✎"
		case "broken":
			icon = "✗"
		case "orphan":
			icon = "⚠"
		}
		detail := ""
		if f.Detail != "" {
			detail = " (" + f.Detail + ")"
		}
		fmt.Fprintf(w, "  %-10s %-40s %-20s %s%s\n", f.Status, f.File, f.Function, icon, detail)
	}

	gen, preserve, broken, orphan := contract.Summary(funcs)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  gen:      %d functions\n", gen)
	fmt.Fprintf(w, "  preserve: %d functions\n", preserve)
	if broken > 0 {
		fmt.Fprintf(w, "  broken:   %d functions\n", broken)
	}
	if orphan > 0 {
		fmt.Fprintf(w, "  orphan:   %d functions\n", orphan)
	}
}
