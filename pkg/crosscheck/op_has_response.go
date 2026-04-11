//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what opHasResponse — operation 목록에서 opID 일치 + 특정 status code 응답 존재 여부 확인
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func opHasResponse(ops map[string]*openapi3.Operation, opID, code string) bool {
	for _, op := range ops {
		if op.OperationID != opID || op.Responses == nil {
			continue
		}
		if op.Responses.Value(code) != nil {
			return true
		}
	}
	return false
}
