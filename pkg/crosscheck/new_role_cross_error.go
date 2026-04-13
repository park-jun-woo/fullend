//ff:func feature=crosscheck type=util control=sequence topic=policy-check
//ff:what newRoleCrossError — X-76 위반 CrossError 공통 생성

package crosscheck

import "fmt"

func newRoleCrossError(funcName, role string) CrossError {
	return CrossError{
		Rule:       "X-76",
		Context:    fmt.Sprintf("%s.ssac Role=%q", funcName, role),
		Level:      "WARNING",
		Message:    fmt.Sprintf("하드코딩 Role %q 는 OPA 정책의 어떤 allow 규칙에서도 명시되지 않음 — role 제약 액션은 모두 차단됨", role),
		Suggestion: "OPA 정책에 해당 role 용 allow 규칙 추가, 또는 SSaC 에서 OPA 가 인정하는 role 로 변경",
	}
}
