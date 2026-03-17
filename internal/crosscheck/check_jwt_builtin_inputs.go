//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what jwtBuiltinFuncs 대상 @call의 input key가 claims 필드와 일치하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/projectconfig"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// CheckJWTBuiltinInputs validates that @call inputs for jwt builtin functions
// use keys that match claims field names.
func CheckJWTBuiltinInputs(serviceFuncs []ssacparser.ServiceFunc, claims map[string]projectconfig.ClaimDef) []CrossError {
	if claims == nil {
		return nil
	}

	claimFields := make(map[string]bool, len(claims))
	for field := range claims {
		claimFields[field] = true
	}

	var errs []CrossError
	for _, sf := range serviceFuncs {
		for seqIdx, seq := range sf.Sequences {
			if seq.Type != "call" || seq.Model == "" {
				continue
			}
			_, _, key := parseCallKey(seq.Model)
			if !jwtBuiltinFuncs[key] {
				continue
			}
			for inputKey := range seq.Inputs {
				if !claimFields[inputKey] {
					errs = append(errs, CrossError{
						Rule:    "SSaC @call → Claims",
						Context: fmt.Sprintf("%s seq[%d] @call %s", sf.Name, seqIdx, seq.Model),
						Message: fmt.Sprintf("@call input key %q가 claims 필드에 없습니다 (유효: %s)", inputKey, claimFieldList(claims)),
						Level:   "ERROR",
					})
				}
			}
		}
	}
	return errs
}

func claimFieldList(claims map[string]projectconfig.ClaimDef) string {
	var keys []string
	for k := range claims {
		keys = append(keys, k)
	}
	return fmt.Sprintf("%v", keys)
}
