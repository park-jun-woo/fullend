//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what HTTP 핸들러 함수를 생성 (분석, 본문, import, 조립)
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// generateHTTPFunc는 HTTP 핸들러 함수를 생성한다.
func (g *GoTarget) generateHTTPFunc(sf ssacparser.ServiceFunc, st *rule.Ground) ([]byte, error) {
	ctx := analyzeHTTPFunc(sf, st, g)
	bodyBuf := buildHTTPFuncBody(sf, st, ctx)

	imports := collectImports(sf, ctx.requestParams, ctx.pathParams, ctx.needsCU, ctx.needsQO)
	imports = filterUsedImports(imports, bodyBuf.String())

	return assembleGoSource(ctx.pkgName, imports, bodyBuf.Bytes())
}
