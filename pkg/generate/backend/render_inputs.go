//ff:func feature=rule type=util control=iteration dimension=1
//ff:what renderInputs — 시퀀스의 Inputs 맵을 Go struct literal로 렌더
package backend

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func renderInputs(seq parsessac.Sequence) string {
	var parts []string
	for k, v := range seq.Inputs {
		parts = append(parts, k+": "+v)
	}
	return strings.Join(parts, ", ")
}
