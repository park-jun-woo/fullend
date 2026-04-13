//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkFuncCoverage — func spec이 SSaC @call에서 사용되는지 검증
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkFuncCoverage(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ProjectFuncSpecs) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("func-coverage")
	graph.Rule(rule.CoverageCheck).With(&rule.CoverageCheckSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-62", Level: "WARNING", Message: "func spec not referenced by any SSaC @call"},
		LookupKey: "SSaC.callRef",
	})

	var errs []CrossError
	for _, sp := range fs.ProjectFuncSpecs {
		key := strings.ToLower(sp.Package + "." + sp.Name)
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", key)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, sp.Package+"."+sp.Name)...)
	}
	return errs
}
