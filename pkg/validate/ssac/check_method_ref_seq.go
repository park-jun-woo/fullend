//ff:func feature=rule type=rule control=sequence
//ff:what checkMethodRefSeq — 단일 시퀀스에서 Model.Method 존재 검증 (S-49 내부)
package ssac

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkMethodRefSeq(ground *rule.Ground, file, funcName string, seqIdx int, seq parsessac.Sequence) []validate.ValidationError {
	if seq.Type != "get" && seq.Type != "post" && seq.Type != "put" && seq.Type != "delete" {
		return nil
	}
	model := extractModelFromSeq(seq)
	if model == "" {
		return nil
	}
	methodKey := "SymbolTable.method." + model
	if _, ok := ground.Lookup[methodKey]; !ok {
		return nil
	}
	method := seq.Model[strings.IndexByte(seq.Model, '.')+1:]
	graph := toulmin.NewGraph("method-ref")
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "S-49", Level: "ERROR", Message: "method " + method + " not found on model " + model},
		LookupKey: methodKey,
	})
	ctx := toulmin.NewContext()
	ctx.Set("ground", ground)
	ctx.Set("claim", method)
	results, _ := graph.Evaluate(ctx)
	return toValidationErrors(results, file, funcName, seqIdx)
}
