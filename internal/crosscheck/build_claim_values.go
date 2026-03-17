//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what fullend.yaml claims에서 claim key 집합을 생성
package crosscheck

import "github.com/park-jun-woo/fullend/internal/projectconfig"

func buildClaimValues(claims map[string]projectconfig.ClaimDef) map[string]bool {
	claimValues := make(map[string]bool)
	for _, def := range claims {
		claimValues[def.Key] = true
	}
	return claimValues
}
