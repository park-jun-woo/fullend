//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what claims 맵의 필드명을 문자열 목록으로 반환
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

func claimFieldList(claims map[string]projectconfig.ClaimDef) string {
	var keys []string
	for k := range claims {
		keys = append(keys, k)
	}
	return fmt.Sprintf("%v", keys)
}
