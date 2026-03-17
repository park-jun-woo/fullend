//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 단일 ServiceFunc의 Go 코드를 생성하는 래퍼
package generator

import (
	"github.com/park-jun-woo/fullend/internal/funcspec"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// GenerateFunc는 단일 ServiceFunc의 Go 코드를 생성한다.
func GenerateFunc(sf parser.ServiceFunc, st *validator.SymbolTable, funcSpecs []funcspec.FuncSpec) ([]byte, error) {
	return (&GoTarget{FuncSpecs: funcSpecs}).GenerateFunc(sf, st)
}
