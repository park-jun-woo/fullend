//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=policy-check
//ff:what 리소스명에서 DDL 테이블명 해석
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// resolveTableName finds the DDL table for a resource name.
func resolveTableName(resource string, st *ssacvalidator.SymbolTable) string {
	snake := pascalToSnake(resource)
	candidates := []string{
		inflection.Plural(strings.ToLower(resource)),
		strings.ToLower(resource),
		inflection.Plural(snake),
		snake,
	}
	for _, c := range candidates {
		if _, ok := st.DDLTables[c]; ok {
			return c
		}
	}
	return ""
}
