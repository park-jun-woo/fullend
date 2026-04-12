//ff:func feature=rule type=generator control=sequence
//ff:what emitPut — @put (UPDATE exec) 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitPut(seq parsessac.Sequence) string {
	return fmt.Sprintf("\tif err := s.%s.%s(ctx, %s); err != nil { return nil, err }\n",
		extractModel(seq.Model), extractMethod(seq.Model), renderArgs(seq))
}
