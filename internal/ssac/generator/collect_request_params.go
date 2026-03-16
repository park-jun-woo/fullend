//ff:func feature=ssac-gen type=generator control=iteration dimension=2 topic=request-params
//ff:what 시퀀스에서 request 소스 파라미터를 수집하고 추출 코드를 생성
package generator

import (
	"github.com/ettle/strcase"
	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

func collectRequestParams(seqs []parser.Sequence, st *validator.SymbolTable, pathParamSet map[string]bool, operationID string) []typedRequestParam {
	rawParams := collectRawRequestParams(seqs, st, pathParamSet)

	if shouldUseJSONBody(seqs, st, rawParams) {
		var rs *validator.RequestSchema
		if st != nil && st.RequestSchemas != nil {
			if schema, ok := st.RequestSchemas[operationID]; ok {
				rs = &schema
			}
		}
		return buildJSONBodyParams(rawParams, rs)
	}

	var params []typedRequestParam
	for _, rp := range rawParams {
		varName := strcase.ToGoCamel(rp.name)
		code := generateExtractCode(varName, rp.name, rp.goType)
		params = append(params, typedRequestParam{
			name:        rp.name,
			goType:      rp.goType,
			extractCode: code,
		})
	}
	return params
}
