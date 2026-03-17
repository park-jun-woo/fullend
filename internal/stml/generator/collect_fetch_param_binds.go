//ff:func feature=stml-gen type=util control=iteration dimension=1
//ff:what FetchBlock의 Params와 중첩 Fetch의 ParamBind를 재귀 수집한다
package generator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

func collectFetchParamBinds(f parser.FetchBlock, params []parser.ParamBind) []parser.ParamBind {
	params = append(params, f.Params...)
	for _, child := range f.NestedFetches {
		params = collectFetchParamBinds(child, params)
	}
	return params
}
