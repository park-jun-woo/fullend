//ff:func feature=orchestrator type=util control=sequence
//ff:what AllSSOTKinds returns all SSOT kinds that fullend manages.

package orchestrator

// AllSSOTKinds returns all SSOT kinds that fullend manages.
func AllSSOTKinds() []SSOTKind {
	return []SSOTKind{
		KindConfig, KindOpenAPI, KindDDL, KindSSaC, KindModel,
		KindSTML, KindStates, KindPolicy, KindScenario, KindFunc,
	}
}
