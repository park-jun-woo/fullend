//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateForbiddenRefs — 금지 참조 검증 (S-31~S-35, S-43, S-47)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateForbiddenRefs(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	// Ensure forbidden sets in ground
	ensureForbiddenSets(ground)

	goGraph := toulmin.NewGraph("go-reserved")
	goGraph.Rule(rule.ForbiddenRef).With(&rule.ForbiddenRefSpec{
		BaseSpec: rule.BaseSpec{Rule: "S-34", Level: "ERROR", Message: "Go reserved word used as variable name"},
		LookupKey: "go.reserved",
	})
	reservedGraph := toulmin.NewGraph("reserved-source")
	reservedGraph.Rule(rule.ForbiddenRef).With(&rule.ForbiddenRefSpec{
		BaseSpec: rule.BaseSpec{Rule: "S-33", Level: "ERROR", Message: "reserved source used as result variable"},
		LookupKey: "ssac.reservedSource",
	})
	dotGraph := toulmin.NewGraph("no-dot-prefix")
	dotGraph.Rule(rule.NameFormat).With(&rule.NameFormatSpec{
		BaseSpec: rule.BaseSpec{Rule: "S-47", Level: "ERROR", Message: "package-prefix @model not allowed"},
		Pattern: "no-dot-prefix",
	})

	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, checkSeqForbidden(goGraph, reservedGraph, dotGraph, ground, fn.FileName, fn.Name, i, seq)...)
	}
	return errs
}
