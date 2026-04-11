//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallFuncNameSeqs — 단일 함수의 시퀀스별 @call 함수명 대문자 시작 검증
package crosscheck

import (
	"strings"
	"unicode"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func checkCallFuncNameSeqs(fnName string, seqs []ssac.Sequence) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		if seq.Type != "call" {
			continue
		}
		idx := strings.IndexByte(seq.Model, '.')
		if idx <= 0 || idx+1 >= len(seq.Model) {
			continue
		}
		funcName := seq.Model[idx+1:]
		if len(funcName) > 0 && unicode.IsLower(rune(funcName[0])) {
			errs = append(errs, CrossError{Rule: "X-38", Context: fnName + "/" + seq.Model, Level: "ERROR",
				Message: "@call function name must start with uppercase"})
		}
	}
	return errs
}
