//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what opHas2xx — operation 목록에서 opID 일치 + 2xx 응답 존재 여부 확인
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func opHas2xx(ops map[string]*openapi3.Operation, opID string) bool {
	for _, op := range ops {
		if op.OperationID != opID || op.Responses == nil {
			continue
		}
		if responsesHave2xx(op.Responses.Map()) {
			return true
		}
	}
	return false
}
