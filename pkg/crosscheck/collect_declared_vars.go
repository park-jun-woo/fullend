//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectDeclaredVars — 시퀀스에서 결과 변수로 선언된 이름 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectDeclaredVars(seqs []ssac.Sequence) map[string]bool {
	vars := make(map[string]bool)
	for _, seq := range seqs {
		if seq.Result != nil && seq.Result.Var != "" {
			vars[seq.Result.Var] = true
		}
	}
	return vars
}
