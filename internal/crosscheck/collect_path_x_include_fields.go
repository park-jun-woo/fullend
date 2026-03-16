//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what 단일 PathItem의 모든 Operation에서 x-include 로컬 컬럼명 수집
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// collectPathXIncludeFields collects x-include local fields from a single path item.
func collectPathXIncludeFields(pi *openapi3.PathItem, result map[string]bool) {
	for _, op := range pi.Operations() {
		if op == nil {
			continue
		}
		collectOpXIncludeFields(op, result)
	}
}
