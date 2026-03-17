//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-ddl
//ff:what SSaC input key가 sqlc 메서드 파라미터명과 대소문자까지 일치하는지 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// CheckInputKeyCase validates that SSaC input keys exactly match sqlc method parameter names (case-sensitive).
func CheckInputKeyCase(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError
	for _, fn := range funcs {
		ctx := fmt.Sprintf("%s:%s", fn.FileName, fn.Name)
		for i, seq := range fn.Sequences {
			errs = append(errs, checkSeqInputKeyCase(ctx, i, seq, st)...)
		}
	}
	return errs
}
