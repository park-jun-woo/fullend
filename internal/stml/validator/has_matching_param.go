//ff:func feature=stml-validate type=util control=iteration dimension=1
//ff:what APISymbol의 parameters에 이름이 일치하는 파라미터가 있는지 확인
package validator

import "strings"

func hasMatchingParam(api APISymbol, name string) bool {
	for _, ap := range api.Parameters {
		if strings.EqualFold(ap.Name, name) {
			return true
		}
	}
	return false
}
