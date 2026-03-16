//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=policy-check
//ff:what Rego allow 규칙이 SSaC authorize 시퀀스에 매칭되는지 검증
package crosscheck

import "fmt"

func checkRegoPairsCoverage(allPairs, ssacPairs map[[2]string]bool) []CrossError {
	var errs []CrossError
	for pair := range allPairs {
		if !ssacPairs[pair] {
			errs = append(errs, CrossError{
				Rule:       "Policy ↔ SSaC",
				Context:    fmt.Sprintf("action=%s resource=%s", pair[0], pair[1]),
				Message:    fmt.Sprintf("Rego allow rule (%s, %s) has no matching SSaC authorize sequence — 미사용 정책 룰", pair[0], pair[1]),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("Add @auth \"%s\" \"%s\" sequence to SSaC", pair[0], pair[1]),
			})
		}
	}
	return errs
}
