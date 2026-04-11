//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkSeqErrStatus — 단일 시퀀스의 ErrStatus → OpenAPI 응답 존재 검증 (X-21)
package crosscheck

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func checkSeqErrStatus(funcName string, seq ssac.Sequence, fs *fullend.Fullstack) []CrossError {
	if seq.ErrStatus <= 0 {
		return nil
	}
	code := strconv.Itoa(seq.ErrStatus)
	if !openAPIHasResponse(fs, funcName, code) {
		return []CrossError{{Rule: "X-21", Context: funcName, Level: "ERROR",
			Message: fmt.Sprintf("ErrStatus %s not defined in OpenAPI responses", code)}}
	}
	return nil
}
