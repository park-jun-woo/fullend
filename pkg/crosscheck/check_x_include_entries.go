//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkXIncludeEntries — x-include allowed 항목의 형식/DDL 검증 (X-4, X-5, X-6)
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkXIncludeEntries(g *rule.Ground, path string, allowed []string) []CrossError {
	var errs []CrossError
	for _, entry := range allowed {
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			errs = append(errs, CrossError{Rule: "X-4", Context: path, Level: "ERROR",
				Message: "x-include format error: " + entry + " (expected FKColumn:RefTable.RefColumn)"})
			continue
		}
		refParts := strings.SplitN(parts[1], ".", 2)
		if len(refParts) != 2 {
			errs = append(errs, CrossError{Rule: "X-4", Context: path, Level: "ERROR",
				Message: "x-include format error: " + entry})
			continue
		}
		refTable := refParts[0]
		if !g.Lookup["DDL.table"][refTable] {
			errs = append(errs, CrossError{Rule: "X-5", Context: path, Level: "ERROR",
				Message: "x-include target table " + refTable + " not found in DDL"})
		}
	}
	return errs
}
