//ff:func feature=ssac-gen type=generator control=sequence topic=query-opts
//ff:what QueryOpts 추출 코드를 버퍼에 출력
package ssac

import (
	"bytes"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func writeQueryOptsCode(buf *bytes.Buffer, needsQO bool, funcName string, st *rule.Ground) {
	if needsQO {
		buf.WriteString(generateQueryOptsCode(funcName, st))
		buf.WriteString("\n")
	}
}
