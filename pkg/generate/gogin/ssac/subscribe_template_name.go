//ff:func feature=ssac-gen type=util control=selection topic=subscribe
//ff:what subscribe 함수 내 시퀀스의 템플릿 이름을 반환
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

// subscribeTemplateName은 subscribe 함수 내 시퀀스의 템플릿 이름을 반환한다.
func subscribeTemplateName(seq ssacparser.Sequence) string {
	switch seq.Type {
	case ssacparser.SeqCall:
		if seq.Result != nil {
			return "sub_call_with_result"
		}
		return "sub_call_no_result"
	case ssacparser.SeqPublish:
		return "sub_publish"
	default:
		return "sub_" + seq.Type
	}
}
