//ff:type feature=stml-parse type=model
//ff:what data-sort 파싱 결과를 나타내는 구조체
package stml

// SortDecl represents a parsed data-sort value.
type SortDecl struct {
	Column    string // default sort column
	Direction string // "asc" or "desc"
}
