//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what ServiceFunc 배열에서 모델 사용 정보를 수집
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func collectModelUsages(funcs []parser.ServiceFunc) []modelUsage {
	var usages []modelUsage
	for _, sf := range funcs {
		usages = append(usages, collectModelUsagesFromFunc(sf)...)
	}
	return usages
}
