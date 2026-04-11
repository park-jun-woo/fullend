//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateRequestRef — request.field → OpenAPI request 스키마 존재 검증 (S-50~S-51)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateRequestRef(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	if fn.Subscribe != nil {
		return nil
	}
	reqKey := "OpenAPI.request." + fn.Name
	if _, ok := ground.Lookup[reqKey]; !ok {
		return nil
	}
	graph := toulmin.NewGraph("request-ref")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "S-50", Level: "ERROR", Message: "request field not in OpenAPI request schema"},
		LookupKey: reqKey,
	})

	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, checkRequestRefSeq(graph, ground, fn.FileName, fn.Name, i, seq)...)
	}
	return errs
}
