//ff:func feature=ssac-gen type=generator control=sequence topic=subscribe
//ff:what subscribe 함수의 본문을 생성
package generator

import (
	"bytes"
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func buildSubscribeFuncBody(sf parser.ServiceFunc, st *validator.SymbolTable, g *GoTarget) bytes.Buffer {
	var bodyBuf bytes.Buffer

	msgType := sf.Subscribe.MessageType
	fmt.Fprintf(&bodyBuf, "func (h *Handler) %s(ctx context.Context, message %s) error {\n", sf.Name, msgType)

	resultTypes, varSources := collectResultInfo(sf.Sequences)
	resolver := &FieldTypeResolver{vars: varSources, st: st, fs: g.FuncSpecs}

	useTx := hasWriteSequence(sf.Sequences)
	if useTx {
		writeSubTxBegin(&bodyBuf)
	}

	writeSubscribeSequences(&bodyBuf, sf, st, resultTypes, resolver, useTx)

	if useTx {
		writeSubTxCommit(&bodyBuf)
	}

	bodyBuf.WriteString("\treturn nil\n")
	bodyBuf.WriteString("}\n")
	return bodyBuf
}
