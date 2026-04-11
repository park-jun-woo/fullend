//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what matchDDLValsToEnum — DDL CHECK 값이 OpenAPI enum에 포함되는지 검증
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func matchDDLValsToEnum(table, col string, ddlVals, oapiEnum []string) []CrossError {
	oapiSet := make(rule.StringSet, len(oapiEnum))
	for _, v := range oapiEnum {
		oapiSet[v] = true
	}
	for _, v := range ddlVals {
		if !oapiSet[v] {
			return []CrossError{{Rule: "X-69", Context: table + "." + col, Level: "WARNING",
				Message: "DDL CHECK value " + v + " not in OpenAPI enum [" + strings.Join(oapiEnum, ", ") + "]"}}
		}
	}
	return nil
}
