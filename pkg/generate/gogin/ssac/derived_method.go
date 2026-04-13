//ff:type feature=ssac-gen type=model topic=interface-derive
//ff:what 파생 인터페이스 메서드 정의를 담는 구조체
package ssac

type derivedMethod struct {
	Name         string
	Params       []derivedParam
	HasQueryOpts bool
	ReturnType   string
}
