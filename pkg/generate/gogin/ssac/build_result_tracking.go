//ff:func feature=ssac-gen type=generator control=sequence topic=template-data
//ff:what result 변수 재선언 여부를 추적하여 ReAssign 플래그 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func buildResultTracking(d *templateData, seq ssacparser.Sequence, declaredVars map[string]bool) {
	if seq.Result == nil {
		return
	}
	if declaredVars[seq.Result.Var] {
		d.ReAssign = true
	}
	declaredVars[seq.Result.Var] = true
}
