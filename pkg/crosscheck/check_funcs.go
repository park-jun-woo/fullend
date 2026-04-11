//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkFuncs — SSaC @call → func spec 존재 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkFuncs(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	graph := toulmin.NewGraph("funcs")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-39", Level: "ERROR", Message: "@call references non-existent func spec"},
		LookupKey: "Func.spec",
	})

	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		for _, ref := range collectCallRefs(fn.Sequences, fn.Name) {
			ctx := toulmin.NewContext()
			ctx.Set("ground", g)
			ctx.Set("claim", ref.key)
			results, _ := graph.Evaluate(ctx)
			errs = append(errs, toErrors(results, ref.context)...)
		}
	}
	return errs
}
