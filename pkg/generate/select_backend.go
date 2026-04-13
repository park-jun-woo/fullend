//ff:func feature=genapi type=util control=sequence
//ff:what selectBackend — 백엔드 구현체를 선택한다
package generate

import "github.com/park-jun-woo/fullend/pkg/generate/gogin"

// Backend is the interface implemented by backend code generators.
type Backend interface {
	// Phase004 stub 경계; 실제 시그니처는 후속 작업에서 확정.
}

// selectBackend returns the backend implementation.
// 수요 확인 전까지는 분기 없음 — fullend.yaml 필드 + 분기 로직 동시 도입은 확장 시점에.
func selectBackend() *gogin.GoGin {
	return &gogin.GoGin{}
}
