//ff:type feature=ssac-gen type=util topic=type-resolve
//ff:what 변수.필드의 Go 타입을 조회하는 리졸버 구조체
package ssac

import (
	"github.com/park-jun-woo/fullend/internal/funcspec"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// FieldTypeResolver는 "변수.필드"의 Go 타입을 조회한다.
type FieldTypeResolver struct {
	vars map[string]varSource
	st   *validator.SymbolTable
	fs   []funcspec.FuncSpec
}
