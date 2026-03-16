//ff:type feature=policy type=model
//ff:what OPA Rego 정책 파싱 결과를 담는 구조체
package policy

// Policy represents parsed OPA Rego policy information.
type Policy struct {
	File       string
	Rules      []AllowRule
	Ownerships []OwnershipMapping
	ClaimsRefs []string // all input.claims.xxx references (deduplicated)
}
