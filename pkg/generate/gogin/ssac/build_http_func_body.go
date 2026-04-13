//ff:func feature=ssac-gen type=generator control=sequence topic=http-handler
//ff:what HTTP 핸들러 함수의 본문을 생성
package ssac

import (
	"bytes"
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func buildHTTPFuncBody(sf ssacparser.ServiceFunc, st *rule.Ground, ctx httpFuncContext) bytes.Buffer {
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
