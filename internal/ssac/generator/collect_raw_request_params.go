//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=request-params
//ff:what 시퀀스에서 request 소스의 필드명과 DDL 타입을 수집
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func collectRawRequestParams(seqs []parser.Sequence, st *validator.SymbolTable, pathParamSet map[string]bool) []rawParam {
	seen := map[string]bool{}
	var params []rawParam

	for _, seq := range seqs {
		params = append(params, collectArgsRequestParams(seq, st, pathParamSet, seen)...)
		params = append(params, collectInputsRequestParams(seq, st, pathParamSet, seen)...)
	}
	return params
}
