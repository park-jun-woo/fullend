//ff:func feature=rule type=generator control=sequence
//ff:what emitAuth — @auth 권한 체크 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitAuth(seq parsessac.Sequence) string {
	status := seq.ErrStatus
	if status == 0 {
		status = 403
	}
	return fmt.Sprintf("\tif _, err := authz.Check(authz.CheckRequest{Action: %q, Resource: %q, UserID: currentUser.ID, Role: currentUser.Role, %s}); err != nil { return nil, gin.H{\"error\": %q, \"status\": %d} }\n",
		seq.Action, seq.Resource, renderInputs(seq), seq.Message, status)
}
