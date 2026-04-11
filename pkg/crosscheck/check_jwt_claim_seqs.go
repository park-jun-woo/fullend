//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkJWTClaimSeqs — 단일 함수의 @call 시퀀스에서 currentUser 필드 → Config claims 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkJWTClaimSeqs(graph *toulmin.Graph, g *rule.Ground, funcName string, seqs []ssac.Sequence) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		if seq.Type != "call" {
			continue
		}
		errs = append(errs, checkJWTClaimArgs(graph, g, funcName, seq.Args)...)
	}
	return errs
}
