//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkHurlMethod — Hurl method → OpenAPI method 존재 검증 (X-36)
package crosscheck

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkHurlMethod(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.HurlEntries) == 0 || fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for _, entry := range fs.HurlEntries {
		lookupKey := "OpenAPI.method." + entry.Path
		graph := toulmin.NewGraph("hurl-method")
		graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
			BaseSpec:  rule.BaseSpec{Rule: "X-36", Level: "ERROR", Message: "Hurl method not defined in OpenAPI for this path"},
			LookupKey: lookupKey,
		})
		errs = append(errs, evalRef(graph, g, strings.ToUpper(entry.Method), fmt.Sprintf("%s:%d", entry.File, entry.Line))...)
	}
	return errs
}
