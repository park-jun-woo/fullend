//ff:func feature=rule type=generator control=sequence
//ff:what emitResponse — @response JSON 응답 코드 생성 (shorthand vs explicit)
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitResponse(seq parsessac.Sequence) string {
	if len(seq.Fields) == 0 {
		return fmt.Sprintf("\treturn %s, nil\n", seq.Target)
	}
	return fmt.Sprintf("\treturn %s, nil\n", renderFieldsAsStruct(seq))
}
