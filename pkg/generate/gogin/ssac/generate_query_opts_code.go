//ff:func feature=ssac-gen type=generator control=sequence topic=query-opts
//ff:what 심볼 테이블의 OpenAPI 확장에서 QueryOpts 설정 코드를 생성
package ssac

import (
	"bytes"

	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func generateQueryOptsCode(funcName string, st *validator.SymbolTable) string {
	if st == nil {
		return "\topts := model.ParseQueryOpts(c, model.QueryOptsConfig{})\n"
	}

	op, hasOp := st.Operations[funcName]
	if !hasOp {
		return "\topts := model.ParseQueryOpts(c, model.QueryOptsConfig{})\n"
	}

	var buf bytes.Buffer
	buf.WriteString("\topts := model.ParseQueryOpts(c, model.QueryOptsConfig{\n")

	writePaginationConfig(&buf, op)
	writeSortConfig(&buf, op)
	writeFilterConfig(&buf, op)

	buf.WriteString("\t})\n")
	return buf.String()
}
