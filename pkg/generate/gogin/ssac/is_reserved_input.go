//ff:func feature=ssac-gen type=util control=sequence topic=args-inputs
//ff:what 입력값이 예약 소스(request, currentUser, 리터럴, query)인지 확인
package ssac

import "strings"

func isReservedInput(val string) bool {
	return strings.HasPrefix(val, "request.") ||
		strings.HasPrefix(val, "currentUser.") ||
		strings.HasPrefix(val, `"`) ||
		val == "query"
}
