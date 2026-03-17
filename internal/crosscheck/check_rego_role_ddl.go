//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what Rego role 값이 DDL CHECK 제약의 허용 값과 일치하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/policy"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// CheckRegoRoleDDL validates that Rego input.claims.role values exist in DDL CHECK constraints.
func CheckRegoRoleDDL(policies []*policy.Policy, st *ssacvalidator.SymbolTable) []CrossError {
	ddlRoleValues := collectDDLRoleValues(st)
	if len(ddlRoleValues) == 0 {
		return nil
	}

	regoRoles := collectRegoRoles(policies)
	var errs []CrossError
	for rv, ctx := range regoRoles {
		if !ddlRoleValues[rv] {
			errs = append(errs, CrossError{
				Rule:       "Rego role → DDL CHECK",
				Context:    ctx,
				Message:    fmt.Sprintf("Rego role %q가 DDL CHECK 제약에 없습니다", rv),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("DDL CHECK (role IN (...))에 %q를 추가하거나 Rego를 수정하세요", rv),
			})
		}
	}
	return errs
}
