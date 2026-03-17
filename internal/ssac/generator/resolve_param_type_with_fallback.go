//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what value 기반 타입 추론 후 string이면 key 기반으로 fallback
package generator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func resolveParamTypeWithFallback(val, key, modelName string, st *validator.SymbolTable) string {
	goType := resolveInputParamType(val, modelName, st)
	if goType == "string" && !strings.HasPrefix(val, `"`) {
		if keyType := resolveKeyParamType(key, modelName, st); keyType != "string" {
			return keyType
		}
	}
	return goType
}
