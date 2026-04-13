//ff:func feature=crosscheck type=rule control=iteration dimension=2 topic=func-check
//ff:what checkEmptyNilable — @empty 대상의 @call 바인딩이 pointer 반환인지 검증 (X-75)

package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// checkEmptyNilable verifies that @empty on a @call-bound variable references a
// funcspec with pointer return type. Otherwise `if x == nil` fails to compile.
// Skips @get/@post bindings — model methods return *T by fullend convention.
func checkEmptyNilable(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	_ = g
	var errs []CrossError
	for _, sf := range fs.ServiceFuncs {
		bindings := map[string]*ssac.Sequence{}
		for i := range sf.Sequences {
			seq := &sf.Sequences[i]
			if seq.Result != nil && seq.Result.Var != "" {
				bindings[seq.Result.Var] = seq
			}
			if seq.Type != ssac.SeqEmpty {
				continue
			}
			bind, ok := bindings[rootVar(seq.Target)]
			if !ok || bind.Type != ssac.SeqCall {
				continue
			}
			fn := findFuncSpecByCall(bind.Model, fs.ProjectFuncSpecs, fs.FullendPkgSpecs)
			if fn == nil || fn.ResponsePointer {
				continue
			}
			errs = append(errs, CrossError{
				Rule:       "X-75",
				Context:    fmt.Sprintf("%s.ssac @empty %s", sf.Name, seq.Target),
				Level:      "ERROR",
				Message:    fmt.Sprintf("@empty %q but funcspec %s returns value type (non-pointer); generated `if %s == nil` will not compile", seq.Target, bind.Model, seq.Target),
				Suggestion: fmt.Sprintf("funcspec %s 의 반환을 *%sResponse 로 변경 (pointer 반환)", bind.Model, fn.Name),
			})
		}
	}
	return errs
}
