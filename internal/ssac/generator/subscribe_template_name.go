//ff:func feature=ssac-gen type=util control=selection topic=subscribe
//ff:what subscribe 함수 내 시퀀스의 템플릿 이름을 반환
package generator

import "github.com/park-jun-woo/fullend/internal/ssac/parser"

// subscribeTemplateName은 subscribe 함수 내 시퀀스의 템플릿 이름을 반환한다.
func subscribeTemplateName(seq parser.Sequence) string {
	switch seq.Type {
	case parser.SeqCall:
		if seq.Result != nil {
			return "sub_call_with_result"
		}
		return "sub_call_no_result"
	case parser.SeqPublish:
		return "sub_publish"
	default:
		return "sub_" + seq.Type
	}
}
