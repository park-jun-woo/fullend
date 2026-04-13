//ff:func feature=ssac-gen type=util control=sequence topic=guard
//ff:what 커스텀 에러 코드가 있으면 변환하고 없으면 기본값 반환
package ssac

func guardErrStatus(code int, defaultStatus string) string {
	if code != 0 {
		return httpStatusConst(code)
	}
	return defaultStatus
}
