//ff:func feature=reporter type=formatter control=iteration dimension=1
//ff:what 검증 보고서 전체를 포맷팅하여 출력한다
package reporter

import "io"

// Print writes the formatted validation report to w.
func Print(w io.Writer, r *Report) {
	for _, step := range r.Steps {
		printStep(w, step)
	}
	printSummary(w, r)
}
