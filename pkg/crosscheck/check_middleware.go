//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkMiddleware — Config middleware ↔ OpenAPI securitySchemes 교차 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkMiddleware(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.Manifest == nil || fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError

	graph := toulmin.NewGraph("middleware")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-50", Level: "ERROR", Message: "Config middleware not in OpenAPI securitySchemes"},
		LookupKey: "OpenAPI.security",
	})
	for _, mw := range fs.Manifest.Backend.Middleware {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", mw)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, mw)...)
	}

	return errs
}
