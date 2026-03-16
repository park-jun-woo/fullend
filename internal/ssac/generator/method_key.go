//ff:type feature=ssac-gen type=model
//ff:what 모델-메서드 쌍을 맵 키로 사용하기 위한 구조체
package generator

type methodKey struct {
	model, method string
}
