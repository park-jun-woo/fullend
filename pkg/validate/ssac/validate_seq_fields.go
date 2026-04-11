//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateSeqFields — 단일 시퀀스의 필수/금지 필드를 toulmin으로 평가
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateSeqFields(file, funcName string, seqIdx int, seq parsessac.Sequence, ground *rule.Ground) []validate.ValidationError {
	specs := fieldRequiredSpecs(seq.Type)
	if len(specs) == 0 {
		return nil
	}
	graph := toulmin.NewGraph("field-required-" + seq.Type)
	for _, spec := range specs {
		graph.Rule(rule.FieldRequired).With(spec)
	}

	claim := buildFieldPresence(seq)
	ctx := toulmin.NewContext()
	ctx.Set("ground", ground)
	ctx.Set("claim", claim)
	results, _ := graph.Evaluate(ctx)
	return toValidationErrors(results, file, funcName, seqIdx)
}
