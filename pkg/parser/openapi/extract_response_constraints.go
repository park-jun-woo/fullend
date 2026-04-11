//ff:func feature=manifest type=parser control=iteration dimension=1
//ff:what ExtractResponseConstraints — operationId별 response 필드 제약조건 추출
package openapi

import "github.com/getkin/kin-openapi/openapi3"

// ExtractResponseConstraints returns field constraints for the 2xx response of each operationId.
func ExtractResponseConstraints(doc *openapi3.T) map[string]map[string]FieldConstraint {
	result := make(map[string]map[string]FieldConstraint)
	if doc == nil || doc.Paths == nil {
		return result
	}
	for _, item := range doc.Paths.Map() {
		extractResponseConstraintsOps(result, item)
	}
	return result
}
