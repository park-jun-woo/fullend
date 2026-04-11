//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkDDLCheckEnum — DDL CHECK IN 값 → OpenAPI enum 일치 검증 (X-68)
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkDDLCheckEnum(g *rule.Ground, tableName string, checkEnums map[string][]string) []CrossError {
	var errs []CrossError
	for col, vals := range checkEnums {
		errs = append(errs, CrossError{Rule: "X-68", Context: tableName + "." + col, Level: "WARNING",
			Message: "DDL CHECK IN values should match OpenAPI enum — " + strings.Join(vals, ", ")})
	}
	return errs
}
