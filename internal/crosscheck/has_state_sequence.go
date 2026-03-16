//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=states
//ff:what SSaC 함수에 @state 시퀀스가 있는지 확인
package crosscheck

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// hasStateSequence checks if a function has a @state sequence.
func hasStateSequence(fn ssacparser.ServiceFunc) bool {
	for _, seq := range fn.Sequences {
		if seq.Type == "state" {
			return true
		}
	}
	return false
}
