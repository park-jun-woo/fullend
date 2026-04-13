//ff:type feature=iface-parse type=model
//ff:what Go 인터페이스 파싱 결과 타입
package iface

// Interface는 specs/model/*.go 에 선언된 Go 인터페이스 하나를 나타낸다.
type Interface struct {
	Name    string   // 인터페이스 이름 (e.g. "UserModel")
	Methods []string // 메서드 이름 순서 (원본 선언 순서 보존)
}
