//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what Arg에서 DDL 테이블을 참조하여 Go 파라미터 타입을 추론
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func resolveArgParamType(a ssacparser.Arg, modelName string, st *validator.SymbolTable) string {
	if a.Literal != "" {
		return "string"
	}

	if resolved := resolveSourceFieldType(a.Source, a.Field, st); resolved != "" {
		return resolved
	}

	snakeName := toSnakeCase(a.Field)

	if goType := lookupModelTableColumn(modelName, snakeName, st); goType != "" {
		return goType
	}

	if goType := lookupRefIDType(a.Field, st); goType != "" {
		return goType
	}

	return lookupAllTablesColumn(snakeName, st)
}
