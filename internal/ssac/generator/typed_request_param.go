//ff:type feature=ssac-gen type=model topic=request-params
//ff:what 타입 정보가 포함된 요청 파라미터 구조체
package generator

type typedRequestParam struct {
	name        string
	goType      string
	extractCode string
}
