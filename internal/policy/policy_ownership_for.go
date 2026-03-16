//ff:func feature=policy type=util control=iteration dimension=1 topic=policy-check
//ff:what 특정 리소스의 소유권 매핑을 찾아 반환한다
package policy

// OwnershipFor returns the ownership mapping for a resource, if any.
func (p *Policy) OwnershipFor(resource string) (OwnershipMapping, bool) {
	for _, o := range p.Ownerships {
		if o.Resource == resource {
			return o, true
		}
	}
	return OwnershipMapping{}, false
}
