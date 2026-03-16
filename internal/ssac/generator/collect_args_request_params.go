//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 시퀀스의 Args에서 request 소스 파라미터를 수집
package generator

import (
	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

func collectArgsRequestParams(seq parser.Sequence, st *validator.SymbolTable, pathParamSet map[string]bool, seen map[string]bool) []rawParam {
	var params []rawParam
	for _, a := range seq.Args {
		if a.Source != "request" || seen[a.Field] || pathParamSet[a.Field] {
			continue
		}
		seen[a.Field] = true
		goType := "string"
		if st != nil {
			goType = lookupDDLType(a.Field, st)
		}
		params = append(params, rawParam{a.Field, goType})
	}
	return params
}
