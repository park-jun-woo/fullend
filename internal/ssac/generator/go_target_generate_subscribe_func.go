//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 큐 구독 핸들러 함수를 생성
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// generateSubscribeFunc는 큐 구독 핸들러 함수를 생성한다.
func (g *GoTarget) generateSubscribeFunc(sf parser.ServiceFunc, st *validator.SymbolTable) ([]byte, error) {
	pkgName := "service"
	if sf.Domain != "" {
		pkgName = sf.Domain
	}

	bodyBuf := buildSubscribeFuncBody(sf, st, g)

	imports := collectSubscribeImports(sf)
	imports = filterUsedImports(imports, bodyBuf.String())

	return assembleGoSource(pkgName, imports, bodyBuf.Bytes())
}
