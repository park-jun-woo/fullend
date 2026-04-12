//ff:func feature=rule type=util control=iteration dimension=1
//ff:what renderFieldsAsStruct — 시퀀스의 Fields 맵을 gin.H{} 렌더
package backend

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func renderFieldsAsStruct(seq parsessac.Sequence) string {
	var parts []string
	for k, v := range seq.Fields {
		parts = append(parts, `"`+k+`": `+v)
	}
	return "gin.H{" + strings.Join(parts, ", ") + "}"
}
