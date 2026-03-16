//ff:type feature=genmodel type=model
//ff:what 메서드 정보 타입 정의
package genmodel

type methodInfo struct {
	Name       string
	HTTPMethod string
	Path       string
	Params     []paramInfo
	ReturnType string // empty = error only
}
