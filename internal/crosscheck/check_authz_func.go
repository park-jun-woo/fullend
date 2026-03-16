//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 단일 SSaC 함수의 @auth 시퀀스 입력 필드를 검증
package crosscheck

import (
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func checkAuthzFunc(fn ssacparser.ServiceFunc, expectedFields map[string]bool) []CrossError {
	var errs []CrossError
	for seqIdx, seq := range fn.Sequences {
		if seq.Type != "auth" {
			continue
		}
		errs = append(errs, checkAuthzSeqInputs(fn.Name, seqIdx, seq, expectedFields)...)
	}
	return errs
}
