//ff:func feature=ssac-gen type=util control=selection
//ff:what HTTP 함수 내 시퀀스의 템플릿 이름을 반환
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

// templateName은 HTTP 함수 내 시퀀스의 템플릿 이름을 반환한다.
func templateName(seq parser.Sequence) string {
	switch seq.Type {
	case parser.SeqResponse:
		if seq.Target != "" {
			return "response_direct"
		}
		return "response"
	case parser.SeqCall:
		if seq.Result != nil {
			return "call_with_result"
		}
		return "call_no_result"
	case parser.SeqPublish:
		return "publish"
	default:
		return seq.Type
	}
}
