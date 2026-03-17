//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=http-handler
//ff:what ServiceFunc를 분석하여 HTTP 함수 생성에 필요한 컨텍스트를 구성
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func analyzeHTTPFunc(sf parser.ServiceFunc, st *validator.SymbolTable, g *GoTarget) httpFuncContext {
	pathParams := getPathParams(sf.Name, st)
	pathParamSet := map[string]bool{}
	for _, pp := range pathParams {
		pathParamSet[pp.Name] = true
	}

	requestParams := collectRequestParams(sf.Sequences, st, pathParamSet, sf.Name)
	needsCU := needsCurrentUser(sf.Sequences)
	needsQO := needsQueryOpts(sf, st)

	pkgName := "service"
	if sf.Domain != "" {
		pkgName = sf.Domain
	}

	resultTypes, varSources := collectResultInfo(sf.Sequences)
	resolver := &FieldTypeResolver{vars: varSources, st: st, fs: g.FuncSpecs}

	return httpFuncContext{
		pathParams:    pathParams,
		pathParamSet:  pathParamSet,
		requestParams: requestParams,
		needsCU:       needsCU,
		needsQO:       needsQO,
		pkgName:       pkgName,
		resolver:      resolver,
		resultTypes:   resultTypes,
		varSources:    varSources,
	}
}
