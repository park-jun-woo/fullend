//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 단일 SSaC 함수의 시퀀스별 @result/@param DDL 매칭을 검증
package crosscheck

import (
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func checkSSaCDDLFunc(fn ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, ctx string, dtoTypes map[string]bool) []CrossError {
	var errs []CrossError
	for i, seq := range fn.Sequences {
		if seq.Type == "call" {
			continue
		}
		if seq.Package != "" {
			continue
		}
		if seq.Result != nil && seq.Result.Type != "" {
			errs = append(errs, checkResultType(seq, st, ctx, i, dtoTypes)...)
		}
		if seq.Model != "" {
			errs = append(errs, checkParamTypes(seq, st, ctx, i)...)
		}
	}
	return errs
}
