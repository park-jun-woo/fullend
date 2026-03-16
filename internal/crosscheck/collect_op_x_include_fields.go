//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=openapi-ddl
//ff:what 단일 Operation의 x-include에서 로컬 컬럼명 추출
package crosscheck

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// collectOpXIncludeFields extracts local column names from a single operation's x-include.
func collectOpXIncludeFields(op *openapi3.Operation, result map[string]bool) {
	raw, ok := op.Extensions["x-include"]
	if !ok {
		return
	}
	var includeExt struct {
		Allowed []string `json:"allowed"`
	}
	if err := unmarshalExt(raw, &includeExt); err != nil {
		return
	}
	for _, spec := range includeExt.Allowed {
		colonIdx := strings.Index(spec, ":")
		if colonIdx > 0 {
			result[spec[:colonIdx]] = true
		}
	}
}
