//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what fullend.yaml claims에서 claim key 집합을 생성
package crosscheck

import "github.com/geul-org/fullend/internal/projectconfig"

func buildClaimValues(claims map[string]projectconfig.ClaimDef) map[string]bool {
	claimValues := make(map[string]bool)
	for _, def := range claims {
		claimValues[def.Key] = true
	}
	return claimValues
}
