//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkOpenAPISSaC — OpenAPI operationId → SSaC 함수 사용 여부 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkOpenAPISSaC(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	graph := toulmin.NewGraph("openapi-ssac")
	graph.Rule(rule.CoverageCheck).With(&rule.CoverageCheckSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-16", Level: "WARNING", Message: "OpenAPI operationId has no matching SSaC function"},
		LookupKey: "SSaC.funcName",
	})

	var errs []CrossError
	for opID := range g.Lookup["OpenAPI.operationId"] {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", opID)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, opID)...)
	}
	return errs
}
