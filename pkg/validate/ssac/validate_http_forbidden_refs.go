//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateHTTPForbiddenRefs — HTTP 함수에서 message 사용 금지 (S-44)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateHTTPForbiddenRefs(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("http-forbidden")
	graph.Rule(rule.ForbiddenRef).With(&rule.ForbiddenRefSpec{
		BaseSpec:  rule.BaseSpec{Rule: "S-44", Level: "ERROR", Message: "HTTP function cannot use message"},
		LookupKey: "http.forbidden",
	})
	localG := copyGroundWithVars(ground, fn)
	localG.Lookup["http.forbidden"] = rule.StringSet{"message": true}

	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, checkHTTPForbiddenSeq(graph, localG, fn.FileName, fn.Name, i, seq)...)
	}
	return errs
}
