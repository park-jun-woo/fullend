//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what adaptPolicies — pkg/parser/rego.Policy → internal/policy.Policy 구조 복사
package gogin

import (
	internalpolicy "github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
)

// adaptPolicies copies pkg rego policies to internal policy.Policy slice.
// Phase006 임시 어댑터. 두 타입은 필드 구조가 동일하므로 단순 복사.
// 장기적으로 generator 내부 함수가 pkg/parser/rego 를 직접 받도록 전환하면 제거.
func adaptPolicies(src []rego.Policy) []*internalpolicy.Policy {
	if len(src) == 0 {
		return nil
	}
	dst := make([]*internalpolicy.Policy, 0, len(src))
	for _, p := range src {
		dst = append(dst, &internalpolicy.Policy{
			File:       p.File,
			Rules:      adaptRules(p.Rules),
			Ownerships: adaptOwnerships(p.Ownerships),
			ClaimsRefs: p.ClaimsRefs,
		})
	}
	return dst
}

func adaptRules(src []rego.AllowRule) []internalpolicy.AllowRule {
	dst := make([]internalpolicy.AllowRule, 0, len(src))
	for _, r := range src {
		dst = append(dst, internalpolicy.AllowRule{
			Resource:  r.Resource,
			Actions:   r.Actions,
			RoleValue: r.RoleValue,
			UsesOwner: r.UsesOwner,
			UsesRole:  r.UsesRole,
		})
	}
	return dst
}

func adaptOwnerships(src []rego.OwnershipMapping) []internalpolicy.OwnershipMapping {
	dst := make([]internalpolicy.OwnershipMapping, 0, len(src))
	for _, o := range src {
		dst = append(dst, internalpolicy.OwnershipMapping{
			Resource:  o.Resource,
			Table:     o.Table,
			Column:    o.Column,
			JoinTable: o.JoinTable,
			JoinFK:    o.JoinFK,
		})
	}
	return dst
}
