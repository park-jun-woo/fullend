//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what OpenAPI Operation에 명시적 2xx 응답 코드가 있는지 확인
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// hasExplicit2xx checks if an operation has an explicit 2xx response code.
func hasExplicit2xx(op *openapi3.Operation) bool {
	for code := range op.Responses.Map() {
		if len(code) == 3 && code[0] == '2' {
			return true
		}
	}
	return false
}
