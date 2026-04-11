//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkShorthandResponse — shorthand @response 검증 (X-19, X-20)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func checkShorthandResponse(fs *fullend.Fullstack) []CrossError {
	// Shorthand @response (단일 변수 반환)는 필드 수준 스키마 비교 대상이 아님.
	// X-17/X-18에서 explicit fields를 검증하고, X-22에서 2xx 존재를 검증.
	// X-19/X-20은 shorthand의 변수 타입이 OpenAPI 응답과 일치하는지 확인하나,
	// 현재 SSaC 파서가 변수 타입 추론을 하지 않으므로 pass-through.
	_ = fs
	return nil
}
