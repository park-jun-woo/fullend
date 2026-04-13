//ff:func feature=crosscheck type=rule control=iteration dimension=3 topic=policy-check
//ff:what checkSSaCRoleVsPolicy — SSaC 하드코딩 Role 이 OPA 에서 언급되지 않으면 WARN (X-76)

package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// checkSSaCRoleVsPolicy detects a hardcoded Role: "X" literal in SSaC sequence
// Inputs/Args that is never referenced by any OPA allow rule's RoleValue.
// Severity: WARNING — user created with that role may have intended-role mismatch.
// Overlaps partially with X-64 (config roles ↔ Rego) but catches SSaC-source cases directly.
func checkSSaCRoleVsPolicy(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	_ = g
	if len(fs.ParsedPolicies) == 0 {
		return nil
	}
	opaRoles, _ := collectOPARoleValues(fs)
	var errs []CrossError
	for _, sf := range fs.ServiceFuncs {
		for i := range sf.Sequences {
			seq := &sf.Sequences[i]
			// @post/@put 등 Inputs map
			for field, val := range seq.Inputs {
				if lit, ok := extractQuotedLiteral(val); ok && field == "Role" && !opaRoles[lit] {
					errs = append(errs, newRoleCrossError(sf.Name, lit))
				}
			}
			// @call Args slice
			for _, arg := range seq.Args {
				if arg.Field == "Role" && arg.Literal != "" && !opaRoles[arg.Literal] {
					errs = append(errs, newRoleCrossError(sf.Name, arg.Literal))
				}
			}
		}
	}
	return errs
}

