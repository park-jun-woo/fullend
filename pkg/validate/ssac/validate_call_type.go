//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateCallType — @call input type ↔ FuncRequest field type 매칭 (S-57)
package ssac

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateCallType(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Type != "call" {
			continue
		}
		idx := strings.IndexByte(seq.Model, '.')
		if idx <= 0 {
			continue
		}
		callFunc := seq.Model[idx+1:]
		lookupKey := "Func.request." + callFunc
		if _, ok := ground.Schemas[lookupKey]; !ok {
			continue
		}
		graph := toulmin.NewGraph("call-type-" + callFunc)
		graph.Rule(rule.TypeMatch).With(&rule.TypeMatchSpec{
			BaseSpec:  rule.BaseSpec{Rule: "S-57", Level: "ERROR", Message: "@call input type mismatch"},
			LookupKey: lookupKey,
		})
		errs = append(errs, evalCallTypeArgs(graph, ground, fn, i, seq, callFunc)...)
	}
	return errs
}
