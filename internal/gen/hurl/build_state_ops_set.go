//ff:func feature=gen-hurl type=util
//ff:what Returns operationIDs that have @state annotations in SSaC.
package hurl

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// buildStateOpsSet returns operationIDs that have @state annotations in SSaC.
func buildStateOpsSet(serviceFuncs []ssacparser.ServiceFunc) map[string]bool {
	ops := make(map[string]bool)
	for _, sf := range serviceFuncs {
		for _, seq := range sf.Sequences {
			if seq.Type == ssacparser.SeqState {
				ops[sf.Name] = true
				break
			}
		}
	}
	return ops
}
