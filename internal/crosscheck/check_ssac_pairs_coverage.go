//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what SSaC authorize 쌍이 Rego allow 규칙에 존재하는지 검증
package crosscheck

import "fmt"

func checkSSaCPairsCoverage(ssacPairs, allPairs map[[2]string]bool) []CrossError {
	var errs []CrossError
	for pair := range ssacPairs {
		if !allPairs[pair] {
			errs = append(errs, CrossError{
				Rule:       "Policy ↔ SSaC",
				Context:    fmt.Sprintf("action=%s resource=%s", pair[0], pair[1]),
				Message:    fmt.Sprintf("SSaC authorize (%s, %s) has no matching Rego allow rule — 런타임에 모든 요청 거부됨", pair[0], pair[1]),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("Add allow rule for action=%q resource=%q in policy/*.rego", pair[0], pair[1]),
			})
		}
	}
	return errs
}
