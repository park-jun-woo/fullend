//ff:type feature=crosscheck type=util
//ff:what Hurl 요청/응답 쌍 타입 정의
package crosscheck

// hurlEntry represents one request/response pair extracted from a .hurl file.
type hurlEntry struct {
	Method     string
	Path       string
	StatusCode string
	File       string
	Line       int
}
