//ff:func feature=rule type=generator control=sequence
//ff:what emitPaginateCursor — cursor pagination @get 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitPaginateCursor(seq parsessac.Sequence) string {
	return fmt.Sprintf("\topts := pagination.ParseCursorOpts(c)\n\t%s, err := s.%s.%s(ctx, opts)\n\tif err != nil { return nil, err }\n",
		seq.Result.Var, extractModel(seq.Model), extractMethod(seq.Model))
}
