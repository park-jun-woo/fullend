//ff:func feature=ssac-gen type=util control=selection
//ff:what 메서드 정보와 사용 정보에서 Go 반환 타입을 파생
package generator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/ssac/validator"
)

func deriveReturnType(mi validator.MethodInfo, usage modelUsage, hasQueryOpts bool) string {
	if usage.Result != nil && usage.Result.Wrapper != "" {
		return fmt.Sprintf("(*pagination.%s[%s], error)", usage.Result.Wrapper, usage.Result.Type)
	}

	switch mi.Cardinality {
	case "exec":
		return "error"
	case "many":
		typeName := extractListTypeName(usage)
		if hasQueryOpts {
			return fmt.Sprintf("([]%s, int, error)", typeName)
		}
		return fmt.Sprintf("([]%s, error)", typeName)
	default:
		typeName := "interface{}"
		if usage.Result != nil {
			typeName = usage.Result.Type
		}
		return fmt.Sprintf("(*%s, error)", typeName)
	}
}
