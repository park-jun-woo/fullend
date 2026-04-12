//ff:func feature=rule type=generator control=sequence
//ff:what emitFKGet — FK 참조 @get (이전 result 기반 조회) 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitFKGet(seq parsessac.Sequence) string {
	return fmt.Sprintf("\t// FK reference lookup\n\t%s, err := s.%s.%s(ctx, %s)\n\tif err != nil { return nil, err }\n",
		seq.Result.Var, extractModel(seq.Model), extractMethod(seq.Model), renderArgs(seq))
}
