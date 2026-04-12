//ff:func feature=rule type=generator control=sequence
//ff:what emitPost — @post (INSERT + scan) 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitPost(seq parsessac.Sequence) string {
	return fmt.Sprintf("\t%s, err := s.%s.%s(ctx, %s)\n\tif err != nil { return nil, err }\n",
		seq.Result.Var, extractModel(seq.Model), extractMethod(seq.Model), renderArgs(seq))
}
