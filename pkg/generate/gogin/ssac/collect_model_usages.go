//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=model-collect
//ff:what ServiceFunc 배열에서 모델 사용 정보를 수집
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectModelUsages(funcs []ssacparser.ServiceFunc) []modelUsage {
	var usages []modelUsage
	for _, sf := range funcs {
		usages = append(usages, collectModelUsagesFromFunc(sf)...)
	}
	return usages
}
