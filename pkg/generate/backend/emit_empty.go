//ff:func feature=rule type=generator control=sequence
//ff:what emitEmpty — @empty 가드 코드 생성 (nil/zero → 404 반환)
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitEmpty(seq parsessac.Sequence) string {
	status := seq.ErrStatus
	if status == 0 {
		status = 404
	}
	return fmt.Sprintf("\tif %s == nil { return nil, gin.H{\"error\": %q, \"status\": %d} }\n",
		seq.Target, seq.Message, status)
}
