//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkAuthzInputSeqs — 시퀀스에서 @auth 타입만 골라 input 필드 검증
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func checkAuthzInputSeqs(funcName string, seqs []ssac.Sequence) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		if seq.Type != "auth" {
			continue
		}
		errs = append(errs, checkAuthzInputFields(funcName, seq.Inputs)...)
	}
	return errs
}
