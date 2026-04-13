//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkDDLCoverage — DDL 테이블이 SSaC에서 참조되는지 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkDDLCoverage(g *rule.Ground) []CrossError {
	if len(g.Tables) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("ddl-coverage")

	w := graph.Rule(rule.CoverageCheck).With(&rule.CoverageCheckSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-55", Level: "WARNING", Message: "DDL table not referenced by any SSaC function"},
		LookupKey: "SSaC.modelRef",
	})
	d := graph.Except(rule.IsArchived)
	d.Attacks(w)

	var errs []CrossError
	for table := range g.Tables {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", table)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, table)...)
	}
	return errs
}
