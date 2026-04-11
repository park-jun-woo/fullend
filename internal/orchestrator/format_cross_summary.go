//ff:func feature=orchestrator type=util control=selection
//ff:what formatCrossSummary — ERROR/WARNING 수로 요약 문자열 생성
package orchestrator

import "fmt"

func formatCrossSummary(errCount, warnCount int) string {
	switch {
	case errCount > 0:
		return fmt.Sprintf("%d errors, %d warnings", errCount, warnCount)
	case warnCount > 0:
		return fmt.Sprintf("%d warnings", warnCount)
	default:
		return "0 mismatches"
	}
}
