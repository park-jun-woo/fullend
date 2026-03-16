//ff:type feature=policy type=model topic=policy-check
//ff:what allow 규칙에서 추출한 액션-리소스 쌍 구조체
package policy

// AllowRule represents an extracted (action, resource) pair from an allow rule.
type AllowRule struct {
	Actions    []string // single or set of actions
	Resource   string
	UsesOwner  bool // references input.resource_owner
	UsesRole   bool // references input.user.role
	RoleValue  string
	SourceLine int
}
