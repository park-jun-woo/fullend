//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what 단일 시퀀스의 @call input key가 claims 필드에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func checkSeqJWTInputs(sfName string, seqIdx int, seq ssacparser.Sequence, claimFields map[string]bool, claims map[string]projectconfig.ClaimDef) []CrossError {
	var errs []CrossError
	for inputKey := range seq.Inputs {
		if !claimFields[inputKey] {
			errs = append(errs, CrossError{
				Rule:    "SSaC @call → Claims",
				Context: fmt.Sprintf("%s seq[%d] @call %s", sfName, seqIdx, seq.Model),
				Message: fmt.Sprintf("@call input key %q가 claims 필드에 없습니다 (유효: %s)", inputKey, claimFieldList(claims)),
				Level:   "ERROR",
			})
		}
	}
	return errs
}
