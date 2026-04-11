//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what guessTableFromOp — operationId에서 CRUD 접두사를 제거하여 DDL 테이블명 추정
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"
)

func guessTableFromOp(opID string) string {
	name := opID
	for _, prefix := range []string{"List", "Get", "Create", "Update", "Delete"} {
		if strings.HasPrefix(name, prefix) {
			name = name[len(prefix):]
			break
		}
	}
	return strings.ToLower(inflection.Plural(name))
}
