//ff:func feature=ssac-gen type=util control=selection topic=http-handler
//ff:what HTTP 함수 내 시퀀스의 템플릿 이름을 반환
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

// templateName은 HTTP 함수 내 시퀀스의 템플릿 이름을 반환한다.
func templateName(seq ssacparser.Sequence) string {
	switch seq.Type {
	case ssacparser.SeqResponse:
		if seq.Target != "" {
			return "response_direct"
		}
		return "response"
	case ssacparser.SeqCall:
		if seq.Result != nil {
			return "call_with_result"
		}
		return "call_no_result"
	case ssacparser.SeqPublish:
		return "publish"
	default:
		return seq.Type
	}
}
