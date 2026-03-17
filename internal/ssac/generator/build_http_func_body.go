//ff:func feature=ssac-gen type=generator control=sequence topic=http-handler
//ff:what HTTP 핸들러 함수의 본문을 생성
package generator

import (
	"bytes"
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func buildHTTPFuncBody(sf parser.ServiceFunc, st *validator.SymbolTable, ctx httpFuncContext) bytes.Buffer {
	var bodyBuf bytes.Buffer

	fmt.Fprintf(&bodyBuf, "func (h *Handler) %s(c *gin.Context) {\n", sf.Name)

	writePathParams(&bodyBuf, ctx.pathParams)
	writeCurrentUser(&bodyBuf, ctx.needsCU)
	writeRequestParamsCode(&bodyBuf, ctx.requestParams)
	writeQueryOptsCode(&bodyBuf, ctx.needsQO, sf.Name, st)

	useTx := hasWriteSequence(sf.Sequences)
	if useTx {
		writeTxBegin(&bodyBuf)
	}

	writeHTTPSequences(&bodyBuf, sf, st, ctx, useTx)

	bodyBuf.WriteString("}\n")
	return bodyBuf
}
