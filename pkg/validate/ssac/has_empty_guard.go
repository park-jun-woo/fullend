//ff:func feature=rule type=util control=iteration dimension=1
//ff:what hasEmptyGuard — 이후 시퀀스에 해당 변수의 @empty 가드가 있는지 확인
package ssac

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func hasEmptyGuard(remaining []parsessac.Sequence, varName string) bool {
	for _, seq := range remaining {
		if seq.Type == "empty" && seq.Target == varName {
			return true
		}
	}
	return false
}
