//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 도메인별 모델 맵의 키를 정렬하여 슬라이스로 변환
package generator

import "sort"

func sortDomainModels(domainSet map[string]map[string]bool) map[string][]string {
	result := map[string][]string{}
	for domain, models := range domainSet {
		keys := make([]string, 0, len(models))
		for k := range models {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		result[domain] = keys
	}
	return result
}
