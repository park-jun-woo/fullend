//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateVariableFlow — 변수 선언 후 사용 검증 (S-27~S-30) + IsImplicitVar defeater
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateVariableFlow(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("var-declared")
	w := graph.Rule(rule.VarDeclared).With(&rule.VarDeclaredSpec{
		BaseSpec: rule.BaseSpec{Rule: "S-27", Level: "ERROR", Message: "variable used before declaration"},
	})
	d := graph.Except(rule.IsImplicitVar)
	d.Attacks(w)

	localG := copyGroundWithVars(ground, fn)
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, evalSeqVarRefs(graph, localG, fn.FileName, fn.Name, i, seq)...)
		trackDeclaredVar(localG, seq)
	}
	return errs
}
