//ff:type feature=gen-hurl type=model
//ff:what orderedStep pairs a scenario step with a sort order key.
package hurl

// orderedStep pairs a scenario step with a sort order key.
type orderedStep struct {
	step  scenarioStep
	order float64
}
