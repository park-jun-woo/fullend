//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=states
//ff:what @state 시퀀스가 있는 함수명 수집
package crosscheck

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

// collectGuardStateFuncs collects function names that have @state sequences.
func collectGuardStateFuncs(funcs []ssacparser.ServiceFunc) map[string]bool {
	result := make(map[string]bool)
	for _, fn := range funcs {
		if hasStateSequence(fn) {
			result[fn.Name] = true
		}
	}
	return result
}
