//ff:type feature=ssac-gen type=model
//ff:what 파생 인터페이스 정의를 담는 구조체
package generator

type derivedInterface struct {
	Name    string
	Methods []derivedMethod
}
