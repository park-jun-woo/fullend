//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 심볼 테이블과 SSaC spec을 교차하여 Model interface를 생성하는 래퍼
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// GenerateModelInterfaces는 심볼 테이블과 SSaC spec을 교차하여 Model interface를 생성한다.
func GenerateModelInterfaces(funcs []ssacparser.ServiceFunc, st *validator.SymbolTable, outDir string) error {
	return DefaultTarget().GenerateModelInterfaces(funcs, st, outDir)
}
