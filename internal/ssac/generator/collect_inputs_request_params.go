//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 시퀀스의 Inputs에서 request. 접두사 파라미터를 수집
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

func collectInputsRequestParams(seq parser.Sequence, st *validator.SymbolTable, pathParamSet map[string]bool, seen map[string]bool) []rawParam {
	var params []rawParam
	for _, val := range seq.Inputs {
		if !strings.HasPrefix(val, "request.") {
			continue
		}
		field := val[len("request."):]
		if seen[field] || pathParamSet[field] {
			continue
		}
		seen[field] = true
		goType := "string"
		if st != nil {
			goType = lookupDDLType(field, st)
		}
		params = append(params, rawParam{field, goType})
	}
	return params
}
