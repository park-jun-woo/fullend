//ff:func feature=rule type=util control=iteration dimension=1
//ff:what renderArgs — 시퀀스의 Args를 Go 함수 호출 인자로 렌더
package backend

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func renderArgs(seq parsessac.Sequence) string {
	var parts []string
	for _, arg := range seq.Args {
		parts = append(parts, renderArg(arg))
	}
	return strings.Join(parts, ", ")
}
