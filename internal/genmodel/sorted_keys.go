//ff:func feature=genmodel type=util control=iteration dimension=1
//ff:what OpenAPI 스키마 맵의 키를 정렬하여 반환한다
package genmodel

import (
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
)

func sortedKeys(m openapi3.Schemas) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
