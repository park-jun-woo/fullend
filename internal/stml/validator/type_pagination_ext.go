//ff:type feature=stml-validate type=model
//ff:what x-pagination 확장 정보를 나타내는 구조체
package validator

// PaginationExt represents x-pagination extension.
type PaginationExt struct {
	Style        string // "offset" or "cursor"
	DefaultLimit int
	MaxLimit     int
}
