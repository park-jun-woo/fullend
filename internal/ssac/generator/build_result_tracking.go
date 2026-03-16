//ff:func feature=ssac-gen type=generator control=sequence
//ff:what result 변수 재선언 여부를 추적하여 ReAssign 플래그 설정
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func buildResultTracking(d *templateData, seq parser.Sequence, declaredVars map[string]bool) {
	if seq.Result == nil {
		return
	}
	if declaredVars[seq.Result.Var] {
		d.ReAssign = true
	}
	declaredVars[seq.Result.Var] = true
}
