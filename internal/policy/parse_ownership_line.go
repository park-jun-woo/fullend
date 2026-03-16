//ff:func feature=policy type=parser control=sequence topic=policy-check
//ff:what 한 줄에서 @ownership 어노테이션을 파싱하여 OwnershipMapping을 반환한다
package policy

// parseOwnershipLine attempts to parse an @ownership annotation from a line.
// Returns the mapping and true if found, zero value and false otherwise.
func parseOwnershipLine(line string) (OwnershipMapping, bool) {
	m := reOwnership.FindStringSubmatch(line)
	if m == nil {
		return OwnershipMapping{}, false
	}
	om := OwnershipMapping{
		Resource: m[1],
		Table:    m[2],
		Column:   m[3],
	}
	if m[4] != "" {
		om.JoinTable = m[4]
		om.JoinFK = m[5]
	}
	return om, true
}
