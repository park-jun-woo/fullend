//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what evalReverseCoverage — target 필드가 source에 사용되는지 WARNING 검사
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/rule"

func evalReverseCoverage(openAPIFields, ssacFields []string, ruleID, context string) []CrossError {
	ssacSet := make(rule.StringSet, len(ssacFields))
	for _, f := range ssacFields {
		ssacSet[f] = true
	}
	var errs []CrossError
	for _, f := range openAPIFields {
		if !ssacSet[f] {
			errs = append(errs, CrossError{Rule: ruleID, Context: context, Level: "WARNING",
				Message: "OpenAPI response field " + f + " not in SSaC @response"})
		}
	}
	return errs
}
