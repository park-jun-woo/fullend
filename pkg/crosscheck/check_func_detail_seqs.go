//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkFuncDetailSeqs — 단일 함수의 @call 시퀀스별 세부 검증 위임
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkFuncDetailSeqs(g *rule.Ground, funcName string, seqs []ssac.Sequence, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		if seq.Type != "call" {
			continue
		}
		errs = append(errs, checkCallDetails(g, funcName, seq, fs)...)
	}
	return errs
}
