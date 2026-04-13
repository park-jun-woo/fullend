//ff:func feature=ssac-gen type=generator control=iteration dimension=2 topic=request-params
//ff:what 시퀀스에서 request 소스 파라미터를 수집하고 추출 코드를 생성
package ssac

import (
	"github.com/ettle/strcase"
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func collectRequestParams(seqs []ssacparser.Sequence, st *rule.Ground, pathParamSet map[string]bool, operationID string) []typedRequestParam {
	rawParams := collectRawRequestParams(seqs, st, pathParamSet)

	if shouldUseJSONBody(seqs, st, operationID, rawParams) {
		var rs *rule.RequestSchemaInfo
		if st != nil && st.ReqSchemas != nil {
			if schema, ok := st.ReqSchemas[operationID]; ok {
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
