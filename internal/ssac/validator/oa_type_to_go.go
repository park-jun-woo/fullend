//ff:func feature=symbol type=util
//ff:what OpenAPI type+format을 Go 타입으로 변환한다
package validator

// oaTypeToGo는 OpenAPI type+format을 Go 타입으로 변환한다.
func oaTypeToGo(oaType, format string) string {
	switch oaType {
	case "integer":
		if format == "int64" {
			return "int64"
		}
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "object", "array":
		return "json.RawMessage"
	default: // string, string+uuid 등
		return "string"
	}
}
