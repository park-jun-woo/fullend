//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkFuncErrStatus — 개별 함수의 ErrStatus/응답 OpenAPI 정합성 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func checkFuncErrStatus(funcName string, seqs []ssac.Sequence, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	hasResponse := false
	for _, seq := range seqs {
		if seq.Type == "response" {
			hasResponse = true
		}
		errs = append(errs, checkSeqErrStatus(funcName, seq, fs)...)
	}
	if hasResponse && !openAPIHas2xx(fs, funcName) {
		errs = append(errs, CrossError{Rule: "X-22", Context: funcName, Level: "ERROR",
			Message: "@response exists but OpenAPI has no 2xx response"})
	}
	return errs
}
