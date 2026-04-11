//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkXIncludeAllowed — x-include allowed 항목의 FK 컬럼 존재 여부 검증
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkXIncludeAllowed(path string, allowed []string, fkCols rule.StringSet) []CrossError {
	var errs []CrossError
	for _, entry := range allowed {
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) == 2 && !fkCols[parts[0]] {
			errs = append(errs, CrossError{Rule: "X-6", Context: path, Level: "WARNING",
				Message: "x-include column " + parts[0] + " has no FK constraint in DDL"})
		}
	}
	return errs
}
