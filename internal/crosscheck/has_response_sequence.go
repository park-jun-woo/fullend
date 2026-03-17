//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what SSaC 함수에 @response 시퀀스가 있는지 확인
package crosscheck

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

// hasResponseSequence checks if a function has a @response sequence.
func hasResponseSequence(fn ssacparser.ServiceFunc) bool {
	for _, seq := range fn.Sequences {
		if seq.Type == "response" {
			return true
		}
	}
	return false
}
