//ff:type feature=stml-validate type=model
//ff:what OpenAPI와 custom.ts에서 추출한 심볼 테이블
package validator

// SymbolTable holds all symbols extracted from OpenAPI and custom.ts.
type SymbolTable struct {
	Operations map[string]APISymbol // operationId → APISymbol
}
