//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateVarTypesSeqs — 단일 함수의 시퀀스에서 변수→타입 매핑 수집
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateVarTypesSeqs(g *rule.Ground, funcName string, seqs []ssac.Sequence) {
	for _, seq := range seqs {
		if seq.Result != nil && seq.Result.Var != "" && seq.Result.Type != "" {
			g.Types["SSaC.var."+funcName+"."+seq.Result.Var] = seq.Result.Type
		}
	}
}
