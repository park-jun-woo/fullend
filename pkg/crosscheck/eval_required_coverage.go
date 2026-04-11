//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what evalRequiredCoverage — OpenAPI required 필드가 SSaC에서 사용되는지 검증
package crosscheck

import (
	oapiparser "github.com/park-jun-woo/fullend/pkg/parser/openapi"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func evalRequiredCoverage(g *rule.Ground, opID string, required []string, fields map[string]oapiparser.FieldConstraint) []CrossError {
	usedFields := make(rule.StringSet)
	for name := range fields {
		usedFields[name] = true
	}
	localG := shallowCopyGround(g)
	localG.Lookup["_used"] = usedFields

	graph := toulmin.NewGraph("required-" + opID)
	graph.Rule(rule.CoverageCheck).With(&rule.CoverageCheckSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-66", Level: "WARNING", Message: "OpenAPI required field not used in SSaC"},
		LookupKey: "_used",
	})

	var errs []CrossError
	for _, field := range required {
		errs = append(errs, evalRef(graph, localG, field, opID+"."+field)...)
	}
	return errs
}
