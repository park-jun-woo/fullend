//ff:func feature=rule type=generator control=sequence
//ff:what emitCall — @call 함수 호출 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitCall(seq parsessac.Sequence) string {
	if seq.Result != nil {
		return fmt.Sprintf("\t%s, err := %s(%s)\n\tif err != nil { return nil, err }\n",
			seq.Result.Var, seq.Model, renderArgs(seq))
	}
	return fmt.Sprintf("\tif _, err := %s(%s); err != nil { return nil, err }\n",
		seq.Model, renderArgs(seq))
}
