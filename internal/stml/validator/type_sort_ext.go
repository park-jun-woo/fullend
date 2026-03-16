//ff:type feature=stml-validate type=model
//ff:what x-sort 확장 정보를 나타내는 구조체
package validator

// SortExt represents x-sort extension.
type SortExt struct {
	Allowed   []string
	Default   string
	Direction string
}
