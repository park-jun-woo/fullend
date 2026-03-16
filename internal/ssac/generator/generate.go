//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what ServiceFunc 배열을 받아 outDir에 Go 파일을 생성하는 진입점
package generator

import (
	"github.com/geul-org/fullend/internal/funcspec"
	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

// Generate는 []ServiceFunc를 받아 outDir에 Go 파일을 생성한다.
func Generate(funcs []parser.ServiceFunc, outDir string, st *validator.SymbolTable, funcSpecs []funcspec.FuncSpec) error {
	return GenerateWith(&GoTarget{FuncSpecs: funcSpecs}, funcs, outDir, st)
}
