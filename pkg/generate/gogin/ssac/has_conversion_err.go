//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=request-params
//ff:what 요청 파라미터에 타입 변환이 필요한 항목이 있는지 확인
package ssac

func hasConversionErr(params []typedRequestParam) bool {
	for _, p := range params {
		if p.goType != "string" && p.goType != "json_body" {
			return true
		}
	}
	return false
}
