//ff:func feature=gen-gogin type=util control=iteration dimension=3
//ff:what finds the 2xx success response code for an operationId in OpenAPI and returns Go http.Status constant

package gogin

import "github.com/getkin/kin-openapi/openapi3"

// resolveSuccessStatus finds the 2xx success response code for an operationId in OpenAPI
// and returns the corresponding Go http.Status constant. Returns "" if not found.
func resolveSuccessStatus(doc *openapi3.T, operationID string) string {
	if doc == nil || doc.Paths == nil {
		return ""
	}
	for _, pi := range doc.Paths.Map() {
		for _, op := range pi.Operations() {
			if op.OperationID != operationID {
				continue
			}
			if op.Responses == nil {
				return ""
			}
			for code := range op.Responses.Map() {
				if len(code) == 3 && code[0] == '2' {
					return httpStatusConst(code)
				}
			}
			return ""
		}
	}
	return ""
}
