//ff:type feature=crosscheck type=util
//ff:what OpenAPI 경로 정규화 라우트 타입 정의
package crosscheck

// apiRoute represents a normalized OpenAPI route for hurl matching.
type apiRoute struct {
	Method    string
	Segments  []string
	Responses map[string]bool
}
