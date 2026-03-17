//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what 단일 ServiceFunc의 JWT builtin @call input key 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/projectconfig"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func checkJWTInputsForFunc(sf ssacparser.ServiceFunc, claimFields map[string]bool, claims map[string]projectconfig.ClaimDef) []CrossError {
	var errs []CrossError
	for seqIdx, seq := range sf.Sequences {
		if seq.Type != "call" || seq.Model == "" {
			continue
		}
		_, _, key := parseCallKey(seq.Model)
		if !jwtBuiltinFuncs[key] {
			continue
		}
		errs = append(errs, checkSeqJWTInputs(sf.Name, seqIdx, seq, claimFields, claims)...)
	}
	return errs
}
