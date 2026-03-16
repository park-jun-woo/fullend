//ff:type feature=stml-validate type=model
//ff:what custom.ts에서 추출한 함수 이름 목록
package validator

// CustomSymbol holds exported function names from a custom.ts file.
type CustomSymbol struct {
	Functions map[string]bool
}
