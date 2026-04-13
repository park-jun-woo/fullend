//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what SSaC input KEY 이름에서 DDL 테이블을 참조하여 Go 타입을 추론 (VALUE 기반 실패 시 fallback)
package ssac

import "github.com/park-jun-woo/fullend/internal/ssac/validator"

// resolveKeyParamType resolves a Go type from the SSaC input KEY name against DDL tables.
// Used as fallback when VALUE-based resolution returns "string".
func resolveKeyParamType(key, modelName string, st *validator.SymbolTable) string {
	snakeName := toSnakeCase(key)

	if goType := lookupModelTableColumn(modelName, snakeName, st); goType != "" {
		return goType
	}

	if goType := lookupRefIDType(key, st); goType != "" {
		return goType
	}

	return lookupAllTablesColumn(snakeName, st)
}
