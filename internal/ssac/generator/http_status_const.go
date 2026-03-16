//ff:func feature=ssac-gen type=util control=selection
//ff:what HTTP 상태 코드 정수를 Go net/http 상수 문자열로 변환
package generator

import "fmt"

func httpStatusConst(code int) string {
	switch code {
	case 400:
		return "http.StatusBadRequest"
	case 401:
		return "http.StatusUnauthorized"
	case 402:
		return "http.StatusPaymentRequired"
	case 403:
		return "http.StatusForbidden"
	case 404:
		return "http.StatusNotFound"
	case 409:
		return "http.StatusConflict"
	case 422:
		return "http.StatusUnprocessableEntity"
	case 429:
		return "http.StatusTooManyRequests"
	case 500:
		return "http.StatusInternalServerError"
	case 502:
		return "http.StatusBadGateway"
	case 503:
		return "http.StatusServiceUnavailable"
	default:
		return fmt.Sprintf("%d", code)
	}
}
