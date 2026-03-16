//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=authz-check
//ff:what @auth 시퀀스의 각 입력 키가 CheckRequest 필드에 있는지 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func checkAuthzSeqInputs(funcName string, seqIdx int, seq ssacparser.Sequence, expectedFields map[string]bool) []CrossError {
	ctx := fmt.Sprintf("%s seq[%d] @auth", funcName, seqIdx)

	var errs []CrossError
	for key := range seq.Inputs {
		if !expectedFields[key] {
			errs = append(errs, CrossError{
				Rule:    "Authz ↔ SSaC",
				Context: ctx,
				Message: fmt.Sprintf("@auth input 필드 %q가 CheckRequest에 없음 (가능: Action, Resource, UserID, Role, ResourceID)", key),
				Level:   "ERROR",
			})
		}
	}
	return errs
}
