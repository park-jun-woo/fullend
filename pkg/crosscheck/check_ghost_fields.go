//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkGhostFields — 단일 operation의 response field가 DDL 컬럼에 있는지 검증
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/rule"

func checkGhostFields(op, table string, fields []string, cols rule.StringSet) []CrossError {
	var errs []CrossError
	for _, f := range fields {
		if !cols[f] && f != "id" {
			errs = append(errs, CrossError{Rule: "X-9", Context: op, Level: "WARNING",
				Message: "OpenAPI property " + f + " not found in DDL table " + table})
		}
	}
	return errs
}
