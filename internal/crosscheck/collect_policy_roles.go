//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what 단일 Rego 정책에서 role 값을 수집
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/policy"
)

func collectPolicyRoles(p *policy.Policy, regoRoles map[string]string) {
	for _, rule := range p.Rules {
		if rule.RoleValue == "" {
			continue
		}
		if _, exists := regoRoles[rule.RoleValue]; !exists {
			regoRoles[rule.RoleValue] = fmt.Sprintf("%s:%d", p.File, rule.SourceLine)
		}
	}
}
