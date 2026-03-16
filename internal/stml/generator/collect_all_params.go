//ff:func feature=stml-gen type=util control=iteration dimension=1
//ff:what 페이지의 모든 ParamBind를 Fetch, Action, Children에서 수집한다
package generator

import "github.com/geul-org/fullend/internal/stml/parser"

// collectAllParams gathers all ParamBinds from the page.
func collectAllParams(page parser.PageSpec) []parser.ParamBind {
	var params []parser.ParamBind
	for _, f := range page.Fetches {
		params = collectFetchParamBinds(f, params)
	}
	for _, a := range page.Actions {
		params = append(params, a.Params...)
	}
	for _, a := range collectAllActions(page.Children) {
		params = append(params, a.Params...)
	}
	return params
}
