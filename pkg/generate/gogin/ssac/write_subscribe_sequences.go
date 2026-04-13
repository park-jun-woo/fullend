//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=subscribe
//ff:what subscribe 시퀀스를 순회하며 각 시퀀스의 코드를 버퍼에 출력
package ssac

import (
	"bytes"
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func writeSubscribeSequences(buf *bytes.Buffer, sf ssacparser.ServiceFunc, st *rule.Ground, resultTypes map[string]string, resolver *FieldTypeResolver, useTx bool) {
	errDeclared := useTx
	declaredVars := map[string]bool{}
	usedVars := collectUsedVars(sf.Sequences)
	for i, seq := range sf.Sequences {
		data := buildTemplateData(seq, &errDeclared, declaredVars, resultTypes, st, sf.Name, useTx, resolver)
		markUnusedResult(&data, seq, usedVars)
		tmplName := subscribeTemplateName(seq)
		var seqBuf bytes.Buffer
		if err := goTemplates.ExecuteTemplate(&seqBuf, tmplName, data); err != nil {
			fmt.Fprintf(buf, "// ERROR: sequence[%d] %s template failed: %v\n", i, seq.Type, err)
			continue
		}
		buf.Write(seqBuf.Bytes())
		buf.WriteString("\n")
	}
}
