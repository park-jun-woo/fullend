//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what x-include에서 로컬 컬럼명 수집
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// collectXIncludeLocalFields collects local column names from x-include across all operations.
func collectXIncludeLocalFields(doc *openapi3.T) map[string]bool {
	result := make(map[string]bool)
	if doc.Paths == nil {
		return result
	}
	for _, pi := range doc.Paths.Map() {
		collectPathXIncludeFields(pi, result)
	}
	return result
}
