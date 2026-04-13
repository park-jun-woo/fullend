//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what 2xx 성공 상태 코드를 반환한다
package hurl

import "github.com/getkin/kin-openapi/openapi3"

// getSuccessHTTPCode returns the numeric 2xx success code string for an operation (e.g. "200", "201", "204").
// Falls back to "200" if no explicit 2xx is found.
func getSuccessHTTPCode(op *openapi3.Operation) string {
	if op.Responses == nil {
		return "200"
	}
	for code := range op.Responses.Map() {
		if len(code) == 3 && code[0] == '2' {
			return code
		}
	}
	return "200"
}
