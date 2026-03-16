//ff:type feature=stml-validate type=model
//ff:what OpenAPI 파라미터를 나타내는 심볼
package validator

// ParamSymbol represents an OpenAPI parameter.
type ParamSymbol struct {
	Name string // parameter name
	In   string // "path" or "query"
}
