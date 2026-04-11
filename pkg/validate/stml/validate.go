//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — STML PageSpec 전체 검증 (TM-1~TM-12)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Validate checks all STML pages against OpenAPI bindings.
func Validate(pages []parsestml.PageSpec, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("stml-ref")
	w := graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-1", Level: "ERROR", Message: "data-fetch/data-action operationId not in OpenAPI"},
		LookupKey: "OpenAPI.operationId",
	})
	d := graph.Except(rule.IsCustomTS)
	d.Attacks(w)

	var errs []validate.ValidationError
	for _, page := range pages {
		errs = append(errs, validatePage(graph, ground, page)...)
	}
	return errs
}
