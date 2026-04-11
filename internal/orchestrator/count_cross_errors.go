//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what countCrossErrors — CrossError에서 ERROR/WARNING 수 집계
package orchestrator

import pkgcross "github.com/park-jun-woo/fullend/pkg/crosscheck"

func countCrossErrors(cerrs []pkgcross.CrossError) (int, int) {
	errCount, warnCount := 0, 0
	for _, ce := range cerrs {
		if ce.Level == "WARNING" {
			warnCount++
		} else {
			errCount++
		}
	}
	return errCount, warnCount
}
