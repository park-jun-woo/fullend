//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=model-collect
//ff:what 도메인별로 사용되는 모델 이름을 수집
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

// collectDomainModels는 도메인별로 사용되는 모델 이름을 수집한다.
func collectDomainModels(funcs []ssacparser.ServiceFunc) map[string][]string {
	domainSet := map[string]map[string]bool{}
	for _, sf := range funcs {
		collectModelsForFunc(sf, domainSet)
	}
	return sortDomainModels(domainSet)
}
