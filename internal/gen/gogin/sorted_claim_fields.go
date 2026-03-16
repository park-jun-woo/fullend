//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what returns claim field names in sorted order

package gogin

import (
	"sort"

	"github.com/geul-org/fullend/internal/projectconfig"
)

// sortedClaimFields returns claim field names in sorted order.
func sortedClaimFields(claims map[string]projectconfig.ClaimDef) []string {
	var fields []string
	for field := range claims {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}
