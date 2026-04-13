//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what returns claim field names in sorted order

package gogin

import (
	"sort"

	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

// sortedClaimFields returns claim field names in sorted order.
func sortedClaimFields(claims map[string]manifest.ClaimDef) []string {
	var fields []string
	for field := range claims {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	return fields
}
