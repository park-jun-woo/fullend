//ff:func feature=ssac-gen type=util control=sequence topic=path-params
//ff:what 심볼 테이블에서 함수명에 해당하는 경로 파라미터를 조회
package generator

import "github.com/geul-org/fullend/internal/ssac/validator"

func getPathParams(funcName string, st *validator.SymbolTable) []validator.PathParam {
	if st == nil {
		return nil
	}
	if op, ok := st.Operations[funcName]; ok {
		return op.PathParams
	}
	return nil
}
