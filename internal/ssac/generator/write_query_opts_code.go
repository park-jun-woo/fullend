//ff:func feature=ssac-gen type=generator control=sequence topic=query-opts
//ff:what QueryOpts 추출 코드를 버퍼에 출력
package generator

import (
	"bytes"

	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func writeQueryOptsCode(buf *bytes.Buffer, needsQO bool, funcName string, st *validator.SymbolTable) {
	if needsQO {
		buf.WriteString(generateQueryOptsCode(funcName, st))
		buf.WriteString("\n")
	}
}
