//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateModelRefs — Model.Method 참조 검증 (S-48, S-49) + 이름 형식 (S-46, S-47)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateModelRefs(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	modelGraph := toulmin.NewGraph("ssac-model-ref")
	modelGraph.Rule(rule.ModelRefExists).With(&rule.ModelRefExistsSpec{
		BaseSpec: rule.BaseSpec{Rule: "S-48", Level: "ERROR", Message: "model not found in symbol table"},
	})
	upperGraph := toulmin.NewGraph("ssac-result-upper")
	upperGraph.Rule(rule.NameFormat).With(&rule.NameFormatSpec{
		BaseSpec: rule.BaseSpec{Rule: "S-46", Level: "ERROR", Message: "result type must start with uppercase"},
		Pattern:  "uppercase-start",
	})

	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		errs = append(errs, evalModelRef(modelGraph, upperGraph, ground, fn.FileName, fn.Name, i, seq)...)
	}
	return errs
}
