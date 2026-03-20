//ff:func feature=crosscheck type=util control=sequence
//ff:what makeExt: JSON 직렬화로 OpenAPI Extension 값 생성 헬퍼
package crosscheck

import "encoding/json"

func makeExt(v any) any {
	b, _ := json.Marshal(v)
	return json.RawMessage(b)
}
