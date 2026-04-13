//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOpenAPIDDL — x-sort/x-filter 컬럼 → DDL 존재 검증 (X-1, X-3)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkOpenAPIDDL(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for _, claim := range collectXSortFilterClaims(fs) {
		graph := toulmin.NewGraph("openapi-ddl-" + claim.ruleID)
		graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
			BaseSpec:  rule.BaseSpec{Rule: claim.ruleID, Level: "ERROR", Message: claim.message},
			LookupKey: claim.lookupKey,
		})
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", claim.col)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, claim.context)...)
	}
	return errs
}
