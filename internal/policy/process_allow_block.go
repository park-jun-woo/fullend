//ff:func feature=policy type=parser control=sequence
//ff:what allow 블록 하나에서 AllowRule을 추출한다
package policy

// processAllowBlock extracts an AllowRule from a single allow block string.
func processAllowBlock(block string) (AllowRule, bool) {
	rule := AllowRule{}

	// Extract action(s).
	if m := reAction.FindStringSubmatch(block); m != nil {
		if m[1] != "" {
			// Single action: input.action == "create"
			rule.Actions = []string{m[1]}
		} else if m[2] != "" {
			// Action set: input.action in {"update", "delete"}
			rule.Actions = parseActionSet(m[2])
		}
	}

	// Extract resource.
	if m := reResource.FindStringSubmatch(block); m != nil {
		rule.Resource = m[1]
	}

	// Check owner reference.
	rule.UsesOwner = reOwnerRef.MatchString(block)

	// Check role reference.
	if m := reRoleRef.FindStringSubmatch(block); m != nil {
		rule.UsesRole = true
		rule.RoleValue = m[1]
	}

	if len(rule.Actions) > 0 && rule.Resource != "" {
		return rule, true
	}
	return rule, false
}
