//ff:func feature=crosscheck type=rule control=iteration dimension=3 topic=func-check
//ff:what checkCallInputsVsFuncspec — @call Inputs 의 변수 타입 ↔ funcspec Request field 타입 호환 (X-77)

package crosscheck

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// checkCallInputsVsFuncspec verifies that each @call Input's bound variable type
// has a compatible basename with the funcspec Request field type.
// Skips cases where the input is a literal or the variable type can't be resolved.
// ERROR on clear basename mismatch (e.g. []Action vs []ActionInput).
func checkCallInputsVsFuncspec(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, sf := range fs.ServiceFuncs {
		for i := range sf.Sequences {
			seq := &sf.Sequences[i]
			if seq.Type != ssac.SeqCall {
				continue
			}
			fn := findFuncSpecByCall(seq.Model, fs.ProjectFuncSpecs, fs.FullendPkgSpecs)
			if fn == nil {
				continue
			}
			for field, val := range seq.Inputs {
				if _, quoted := extractQuotedLiteral(val); quoted {
					continue // literal OK
				}
				if strings.Contains(val, ".") {
					continue // field access (x.Y) — field type resolution 미구현, 보수적 skip
				}
				varType := g.Types["SSaC.var."+sf.Name+"."+val]
				if varType == "" {
					continue // unresolved
				}
				reqFieldType := findRequestFieldType(fn, field)
				if reqFieldType == "" {
					continue
				}
				if extractBareTypeName(varType) == extractBareTypeName(reqFieldType) {
					continue
				}
				errs = append(errs, CrossError{
					Rule:       "X-77",
					Context:    fmt.Sprintf("%s.ssac @call %s.%s=%s", sf.Name, seq.Model, field, val),
					Level:      "ERROR",
					Message:    fmt.Sprintf("@call 인자 타입 %s 가 funcspec %s.Request.%s 타입 %s 와 호환 불가", varType, seq.Model, field, reqFieldType),
					Suggestion: fmt.Sprintf("funcspec %s 의 Request.%s 타입을 %s 로 맞추거나 SSaC 에서 변환 단계 추가", seq.Model, field, varType),
				})
			}
		}
	}
	return errs
}
