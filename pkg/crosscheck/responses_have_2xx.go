//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what responsesHave2xx — 응답 코드 맵에 2xx가 있는지 확인
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func responsesHave2xx(responses map[string]*openapi3.ResponseRef) bool {
	for code := range responses {
		if len(code) > 0 && code[0] == '2' {
			return true
		}
	}
	return false
}
