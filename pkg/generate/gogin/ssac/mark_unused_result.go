//ff:func feature=ssac-gen type=util control=sequence topic=template-data
//ff:what 사용되지 않는 result 변수에 Unused/ReAssign 플래그를 설정
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func markUnusedResult(d *templateData, seq ssacparser.Sequence, usedVars map[string]bool) {
	if seq.Result == nil {
		return
	}
	if usedVars[seq.Result.Var] {
		return
	}
	d.Unused = true
	if d.ErrDeclared {
		d.ReAssign = true
	}
}
