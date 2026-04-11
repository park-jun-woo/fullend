//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOwnershipVia — @ownership via join table/FK → DDL 존재 검증 (X-33)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkOwnershipVia(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ParsedPolicies) == 0 || len(fs.DDLTables) == 0 {
		return nil
	}
	tableGraph := toulmin.NewGraph("ownership-via-table")
	tableGraph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-33", Level: "ERROR", Message: "@ownership via join table not found in DDL"},
		LookupKey: "DDL.table",
	})
	var errs []CrossError
	for _, p := range fs.ParsedPolicies {
		errs = append(errs, checkOwnershipViaJoins(tableGraph, g, p.Ownerships)...)
	}
	return errs
}
