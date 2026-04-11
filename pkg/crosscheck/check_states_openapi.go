//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkStatesOpenAPI — States transition → OpenAPI operationId 존재 검증 (X-25)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkStatesOpenAPI(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.StateDiagrams) == 0 || fs.OpenAPIDoc == nil {
		return nil
	}
	graph := toulmin.NewGraph("states-openapi")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-25", Level: "ERROR", Message: "States transition event has no matching OpenAPI operationId"},
		LookupKey: "OpenAPI.operationId",
	})
	var errs []CrossError
	for _, sd := range fs.StateDiagrams {
		for _, tr := range sd.Transitions {
			errs = append(errs, evalRef(graph, g, tr.Event, sd.ID+"."+tr.Event)...)
		}
	}
	return errs
}
