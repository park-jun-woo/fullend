//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 도메인별 Handler struct를 생성하는 래퍼
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// GenerateHandlerStruct는 도메인별 Handler struct를 생성한다.
func GenerateHandlerStruct(funcs []parser.ServiceFunc, st *validator.SymbolTable, outDir string) error {
	return DefaultTarget().GenerateHandlerStruct(funcs, st, outDir)
}
