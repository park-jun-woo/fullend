//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what SSaC 함수의 @response 필드 키 추출
package crosscheck

import (
	"sort"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// extractResponseFieldKeys returns the @response field keys for a function,
// or nil if the function uses shorthand (@response varName) or has no @response.
func extractResponseFieldKeys(fn ssacparser.ServiceFunc) []string {
	for _, seq := range fn.Sequences {
		if seq.Type != "response" {
			continue
		}
		if seq.Target != "" {
			return nil
		}
		if len(seq.Fields) == 0 {
			return nil
		}
		keys := make([]string, 0, len(seq.Fields))
		for k := range seq.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}
	return nil
}
