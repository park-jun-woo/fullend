package crosscheck

// Run executes all cross-validation rules and returns collected errors.
func Run(input *CrossValidateInput) []CrossError {
	return RunRules(input, nil)
}

// RunRules executes rules, skipping names in skipRules.
func RunRules(input *CrossValidateInput, skipRules map[string]bool) []CrossError {
	var errs []CrossError
	for _, r := range rules {
		if skipRules[r.Name] {
			continue
		}
		if r.Requires(input) {
			errs = append(errs, r.Check(input)...)
		}
	}
	return errs
}

// Rules returns the registered rule list (for status/reporting).
func Rules() []Rule {
	return rules
}
