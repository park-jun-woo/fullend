//ff:type feature=ssac-gen type=model topic=interface-derive
//ff:what 모델-메서드 쌍을 맵 키로 사용하기 위한 구조체
package ssac

type methodKey struct {
	model, method string
}
