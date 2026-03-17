//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=http-handler
//ff:what HTTP 시퀀스를 순회하며 각 시퀀스의 코드를 버퍼에 출력
package generator

import (
	"bytes"
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func writeHTTPSequences(buf *bytes.Buffer, sf parser.ServiceFunc, st *validator.SymbolTable, ctx httpFuncContext, useTx bool) {
	errDeclared := hasConversionErr(ctx.requestParams)
	if useTx {
		errDeclared = true
	}
	declaredVars := map[string]bool{}
	funcHasTotal := false
	usedVars := collectUsedVars(sf.Sequences)
	committed := false

	for i, seq := range sf.Sequences {
		if useTx && seq.Type == parser.SeqResponse && !committed {
			writeTxCommit(buf)
			committed = true
		}
		data := buildTemplateData(seq, &errDeclared, declaredVars, ctx.resultTypes, st, sf.Name, useTx, ctx.resolver)
		if data.HasTotal {
			funcHasTotal = true
		}
		if seq.Type == parser.SeqResponse {
			data.HasTotal = funcHasTotal
		}
		markUnusedResult(&data, seq, usedVars)

		tmplName := templateName(seq)
		var seqBuf bytes.Buffer
		if err := goTemplates.ExecuteTemplate(&seqBuf, tmplName, data); err != nil {
			fmt.Fprintf(buf, "// ERROR: sequence[%d] %s template failed: %v\n", i, seq.Type, err)
			continue
		}
		buf.Write(seqBuf.Bytes())
		buf.WriteString("\n")
	}

	if useTx && !committed {
		writeTxCommit(buf)
	}
}
