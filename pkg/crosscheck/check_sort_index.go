//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkSortIndex — x-sort column에 인덱스 존재 여부 WARNING (X-2)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkSortIndex(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for _, claim := range collectXSortFilterClaims(fs) {
		if claim.ruleID != "X-1" {
			continue
		}
		idxKey := "DDL.index." + claim.lookupKey[len("DDL.column."):]
		indexed := g.Lookup[idxKey]
		if !indexed[claim.col] {
			errs = append(errs, CrossError{Rule: "X-2", Context: claim.context, Level: "WARNING",
				Message: "x-sort column " + claim.col + " has no index"})
		}
	}
	return errs
}
