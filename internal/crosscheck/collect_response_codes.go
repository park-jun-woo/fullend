//ff:func feature=crosscheck type=util control=sequence topic=ssac-openapi
//ff:what Operation에서 응답 코드 수집
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// collectResponseCodes collects response codes from an operation.
func collectResponseCodes(op *openapi3.Operation) map[string]bool {
	codes := make(map[string]bool)
	if op.Responses != nil {
		for code := range op.Responses.Map() {
			codes[code] = true
		}
	}
	return codes
}
