//ff:func feature=rule type=util control=sequence
//ff:what trackDeclaredVar — 시퀀스 결과 변수를 Ground.Vars에 추가
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func trackDeclaredVar(g *rule.Ground, seq parsessac.Sequence) {
	if seq.Result != nil && seq.Result.Var != "" {
		g.Vars[seq.Result.Var] = true
	}
}
