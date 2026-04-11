//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkHurl — Hurl path/method → OpenAPI 존재 검증
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkHurl(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.HurlEntries) == 0 || fs.OpenAPIDoc == nil {
		return nil
	}
	graph := toulmin.NewGraph("hurl")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-35", Level: "ERROR", Message: "Hurl path not defined in OpenAPI"},
		LookupKey: "OpenAPI.path",
	})

	var errs []CrossError
	for _, entry := range fs.HurlEntries {
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", entry.Path)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, fmt.Sprintf("%s:%d", entry.File, entry.Line))...)
	}
	return errs
}
