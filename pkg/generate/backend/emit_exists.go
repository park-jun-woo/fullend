//ff:func feature=rule type=generator control=sequence
//ff:what emitExists — @exists 가드 코드 생성 (not nil → 409 반환)
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitExists(seq parsessac.Sequence) string {
	status := seq.ErrStatus
	if status == 0 {
		status = 409
	}
	return fmt.Sprintf("\tif %s != nil { return nil, gin.H{\"error\": %q, \"status\": %d} }\n",
		seq.Target, seq.Message, status)
}
