//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 심볼 테이블과 SSaC spec을 교차하여 Model interface를 생성하는 래퍼
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// GenerateModelInterfaces는 심볼 테이블과 SSaC spec을 교차하여 Model interface를 생성한다.
func GenerateModelInterfaces(funcs []parser.ServiceFunc, st *validator.SymbolTable, outDir string) error {
	return DefaultTarget().GenerateModelInterfaces(funcs, st, outDir)
}
