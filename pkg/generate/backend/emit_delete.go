//ff:func feature=rule type=generator control=sequence
//ff:what emitDelete — @delete (DELETE exec) 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitDelete(seq parsessac.Sequence) string {
	return fmt.Sprintf("\tif err := s.%s.%s(ctx, %s); err != nil { return nil, err }\n",
		extractModel(seq.Model), extractMethod(seq.Model), renderArgs(seq))
}
