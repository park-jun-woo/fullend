//ff:func feature=rule type=generator control=sequence
//ff:what emitState — @state 전이 가드 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitState(seq parsessac.Sequence) string {
	status := seq.ErrStatus
	if status == 0 {
		status = 409
	}
	return fmt.Sprintf("\tif !statemachine.%s.CanTransition(%q, %s) { return nil, gin.H{\"error\": %q, \"status\": %d} }\n",
		seq.DiagramID, seq.Transition, renderInputs(seq), seq.Message, status)
}
