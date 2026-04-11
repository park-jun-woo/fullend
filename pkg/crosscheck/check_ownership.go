//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOwnership — @ownership table/column → DDL 존재 검증 (X-31~X-34)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkOwnership(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ParsedPolicies) == 0 || len(fs.DDLTables) == 0 {
		return nil
	}
	tableGraph := toulmin.NewGraph("ownership-table")
	tableGraph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-31", Level: "ERROR", Message: "@ownership table not found in DDL"},
		LookupKey: "DDL.table",
	})

	var errs []CrossError
	for _, p := range fs.ParsedPolicies {
		errs = append(errs, checkOwnershipMappings(tableGraph, g, p.Ownerships)...)
	}
	return errs
}
