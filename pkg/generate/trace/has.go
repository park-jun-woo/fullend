//ff:func feature=rule type=util control=sequence
//ff:what Has — Pattern에서 warrant 활성화 여부 확인
package trace

// Has returns true if the warrant was activated in the pattern.
func (p Pattern) Has(name string) bool {
	return p[name]
}
