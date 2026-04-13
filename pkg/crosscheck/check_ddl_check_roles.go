//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkDDLCheckRoles — Rego role 값이 DDL CHECK 제약에 있는지 검증 (X-65)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkDDLCheckRoles(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.ParsedPolicies) == 0 || len(fs.DDLTables) == 0 {
		return nil
	}
	roleCheckSets := collectDDLRoleCheckSets(fs)
	if len(roleCheckSets) == 0 {
		return nil
	}
	var errs []CrossError
	for rv := range g.Lookup["Rego.roles"] {
		if !roleInAnySets(rv, roleCheckSets) {
			errs = append(errs, CrossError{Rule: "X-65", Context: rv, Level: "WARNING",
				Message: "Rego role " + rv + " not found in DDL CHECK constraint"})
		}
	}
	return errs
}
