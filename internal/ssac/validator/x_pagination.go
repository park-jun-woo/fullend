//ff:type feature=symbol type=model
//ff:what x-pagination 확장
package validator

// XPagination은 x-pagination 확장이다.
type XPagination struct {
	Style        string `yaml:"style"`
	DefaultLimit int    `yaml:"defaultLimit"`
	MaxLimit     int    `yaml:"maxLimit"`
}
