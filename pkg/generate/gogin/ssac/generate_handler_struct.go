//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 도메인별 Handler struct를 생성하는 래퍼
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// GenerateHandlerStruct는 도메인별 Handler struct를 생성한다.
func GenerateHandlerStruct(funcs []ssacparser.ServiceFunc, st *validator.SymbolTable, outDir string) error {
	return DefaultTarget().GenerateHandlerStruct(funcs, st, outDir)
}
