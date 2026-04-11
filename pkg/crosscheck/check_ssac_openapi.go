//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkSSaCOpenAPI — SSaC funcName ↔ OpenAPI operationId 교차 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkSSaCOpenAPI(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	graph := toulmin.NewGraph("ssac-openapi")

	w := graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-15", Level: "ERROR", Message: "SSaC function has no matching OpenAPI operationId"},
		LookupKey: "OpenAPI.operationId",
	})
	d := graph.Except(rule.IsSubscribe)
	d.Attacks(w)

	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		if fn.Subscribe != nil {
			g.Flags["subscribe"] = true
		} else {
			delete(g.Flags, "subscribe")
		}
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", fn.Name)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, fn.Name)...)
	}
	return errs
}
