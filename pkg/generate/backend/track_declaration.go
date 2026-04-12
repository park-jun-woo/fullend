//ff:func feature=rule type=util control=sequence
//ff:what trackDeclaration — 시퀀스 result 변수를 declared 맵에 등록
package backend

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func trackDeclaration(seq parsessac.Sequence, declared map[string]string) {
	if seq.Result != nil && seq.Result.Var != "" {
		declared[seq.Result.Var] = seq.Result.Type
	}
}
