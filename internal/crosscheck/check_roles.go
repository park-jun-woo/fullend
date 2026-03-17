//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what OPA Rego input.role 값이 fullend.yaml auth.roles와 일치하는지 검증
package crosscheck

import (
	"fmt"
	"sort"

	"github.com/park-jun-woo/fullend/internal/policy"
)

// CheckRoles validates that OPA Rego input.role values match fullend.yaml auth.roles.
func CheckRoles(policies []*policy.Policy, roles []string) []CrossError {
	if len(roles) == 0 {
		return nil
	}

	var errs []CrossError

	roleSet := make(map[string]bool, len(roles))
	for _, r := range roles {
		roleSet[r] = true
	}

	regoRoles := collectRegoRoles(policies)

	// Rego role → fullend.yaml roles (ERROR if missing).
	for rv, ctx := range regoRoles {
		if !roleSet[rv] {
			errs = append(errs, CrossError{
				Rule:       "Rego role → fullend.yaml",
				Context:    ctx,
				Message:    fmt.Sprintf("Rego role %q가 fullend.yaml auth.roles에 없습니다", rv),
				Suggestion: fmt.Sprintf("fullend.yaml auth.roles에 %q를 추가하거나 Rego를 수정하세요", rv),
			})
		}
	}

	// fullend.yaml roles → Rego (WARNING if unused).
	var unused []string
	for _, r := range roles {
		if _, used := regoRoles[r]; !used {
			unused = append(unused, r)
		}
	}
	sort.Strings(unused)
	for _, r := range unused {
		errs = append(errs, CrossError{
			Rule:       "fullend.yaml → Rego role",
			Context:    "fullend.yaml",
			Message:    fmt.Sprintf("fullend.yaml auth.roles의 %q가 Rego에서 사용되지 않습니다", r),
			Level:      "WARNING",
			Suggestion: fmt.Sprintf("Rego에 input.role == %q 조건을 추가하거나 roles에서 제거하세요", r),
		})
	}

	return errs
}
