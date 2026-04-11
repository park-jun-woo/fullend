//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallTypeMatchSeqs — 단일 함수의 시퀀스별 @call arg 타입 비교 위임
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkCallTypeMatchSeqs(g *rule.Ground, funcName string, seqs []ssac.Sequence) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		if seq.Type != "call" {
			continue
		}
		errs = append(errs, checkCallArgTypes(g, funcName, seq.Model, seq.Args)...)
	}
	return errs
}
