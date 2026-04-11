//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateSubForbiddenRefs — @subscribe 함수에서 request/query 사용 금지 (S-42~S-43)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateSubForbiddenRefs(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("sub-forbidden")
	graph.Rule(rule.ForbiddenRef).With(&rule.ForbiddenRefSpec{
		BaseSpec:  rule.BaseSpec{Rule: "S-42", Level: "ERROR", Message: "@subscribe cannot use request/query"},
		LookupKey: "subscribe.forbidden",
	})
	localG := copyGroundWithVars(ground, fn)
	localG.Lookup["subscribe.forbidden"] = rule.StringSet{"request": true, "query": true}

	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, checkSubForbiddenSeq(graph, localG, fn.FileName, fn.Name, i, seq)...)
	}
	return errs
}
