//ff:func feature=crosscheck type=util control=sequence
//ff:what SSaC @call model 문자열에서 패키지, camelCase 함수명, 키 추출
package crosscheck

import (
	"strings"

	"github.com/ettle/strcase"
)

// parseCallKey extracts package, camelCase function name, and lookup key from a call model string.
func parseCallKey(model string) (pkg, camelName, key string) {
	callParts := strings.SplitN(model, ".", 2)
	funcName := model
	if len(callParts) == 2 {
		pkg = callParts[0]
		funcName = callParts[1]
	}
	camelName = strcase.ToGoCamel(funcName)
	key = pkg + "." + camelName
	if pkg == "" {
		key = camelName
	}
	return pkg, camelName, key
}
