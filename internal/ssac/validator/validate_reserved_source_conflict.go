//ff:func feature=ssac-validate type=rule control=iteration dimension=1 topic=string-convert
//ff:what result 변수명과 예약 소스 충돌 검증

package validator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// reservedSources는 사용자가 result 변수명으로 사용할 수 없는 예약 소스다.
var reservedSources = map[string]bool{
	"request":     true,
	"currentUser": true,
	"config":      true,
	"query":       true,
	"message":     true,
}

// validateReservedSourceConflict는 result 변수명이 예약 소스와 충돌하면 ERROR를 반환한다.
func validateReservedSourceConflict(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	for i, seq := range sf.Sequences {
		if seq.Result == nil {
			continue
		}
		if reservedSources[seq.Result.Var] {
			ctx := errCtx{sf.FileName, sf.Name, i}
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q는 예약 소스이므로 result 변수명으로 사용할 수 없습니다", seq.Result.Var)))
		}
	}
	return errs
}
