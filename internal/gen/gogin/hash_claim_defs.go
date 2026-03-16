//ff:func feature=gen-gogin type=util control=iteration
//ff:what computes a contract hash for ClaimDef claims

package gogin

import (
	"strings"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/projectconfig"
)

// HashClaimDefs computes a contract hash for ClaimDef claims.
func HashClaimDefs(claims map[string]projectconfig.ClaimDef) string {
	fields := sortedClaimFields(claims)
	var parts []string
	for _, f := range fields {
		def := claims[f]
		parts = append(parts, f+":"+def.Key+":"+def.GoType)
	}
	return contract.Hash7(strings.Join(parts, ","))
}
