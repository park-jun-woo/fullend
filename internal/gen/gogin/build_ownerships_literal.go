//ff:func feature=gen-gogin type=generator control=iteration
//ff:what generates Go code for []authz.OwnershipMapping literal

package gogin

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/policy"
)

// buildOwnershipsLiteral generates Go code for []authz.OwnershipMapping literal.
func buildOwnershipsLiteral(policies []*policy.Policy) string {
	var mappings []string
	for _, p := range policies {
		for _, om := range p.Ownerships {
			mappings = append(mappings, fmt.Sprintf(
				`{Resource: %q, Table: %q, Column: %q}`,
				om.Resource, om.Table, om.Column,
			))
		}
	}
	if len(mappings) == 0 {
		return "nil"
	}
	return "[]authz.OwnershipMapping{\n\t\t" + strings.Join(mappings, ",\n\t\t") + ",\n\t}"
}
