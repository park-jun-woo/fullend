//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what fullend/project func spec을 key 맵으로 변환
package crosscheck

import "github.com/geul-org/fullend/internal/funcspec"

// buildFuncSpecMap builds a lookup map: "package.funcName" -> FuncSpec.
// Project custom overrides fullend default.
func buildFuncSpecMap(fullendPkgSpecs, projectFuncSpecs []funcspec.FuncSpec) map[string]*funcspec.FuncSpec {
	specMap := make(map[string]*funcspec.FuncSpec)
	for i := range fullendPkgSpecs {
		key := fullendPkgSpecs[i].Package + "." + fullendPkgSpecs[i].Name
		specMap[key] = &fullendPkgSpecs[i]
	}
	for i := range projectFuncSpecs {
		key := projectFuncSpecs[i].Package + "." + projectFuncSpecs[i].Name
		specMap[key] = &projectFuncSpecs[i]
	}
	return specMap
}
