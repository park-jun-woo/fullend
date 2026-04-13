//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=interface-derive
//ff:what 메서드 배열의 반환 타입에 특정 문자열이 포함되어 있는지 확인
package ssac

import "strings"

func methodsHaveReturnSubstring(methods []derivedMethod, substr string) bool {
	for _, m := range methods {
		if strings.Contains(m.ReturnType, substr) {
			return true
		}
	}
	return false
}
