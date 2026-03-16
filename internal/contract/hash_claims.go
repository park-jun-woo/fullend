//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 미들웨어 클레임의 계약 해시를 계산한다
package contract

import (
	"sort"
	"strings"

	"github.com/geul-org/fullend/internal/projectconfig"
)

// HashClaims computes a contract hash for middleware claims (CurrentUser).
// Based on: sorted field:key:type triples.
func HashClaims(claims map[string]projectconfig.ClaimDef) string {
	keys := make([]string, 0, len(claims))
	for k := range claims {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		def := claims[k]
		parts = append(parts, k+":"+def.Key+":"+def.GoType)
	}
	return Hash7(strings.Join(parts, ","))
}
