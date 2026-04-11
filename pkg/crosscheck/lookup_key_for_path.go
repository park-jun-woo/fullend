//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what lookupKeyForPath — operationId에서 모델명을 추론하여 DDL.column 키 생성
package crosscheck

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jinzhu/inflection"
)

func lookupKeyForPath(op *openapi3.Operation) string {
	if op.OperationID == "" {
		return ""
	}
	name := op.OperationID
	for _, prefix := range []string{"List", "Get", "Create", "Update", "Delete"} {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
			break
		}
	}
	table := strings.ToLower(inflection.Plural(name))
	return "DDL.column." + table
}
