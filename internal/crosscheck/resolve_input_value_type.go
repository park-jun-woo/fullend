//ff:func feature=crosscheck type=util control=selection topic=func-check
//ff:what Input 값 문자열의 Go 타입을 해석
package crosscheck

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// resolveInputValueType resolves the Go type of an Input value string.
// Patterns: "request.Field" -> OpenAPI, "var.Field" -> DDL/FuncResponse, "\"literal\"" -> string,
// bare variable -> definedVars type, numeric/bool/nil literal -> inferred type, "currentUser.*" -> skip.
func resolveInputValueType(value string, definedVars map[string]string, st *ssacvalidator.SymbolTable, doc *openapi3.T, funcName string, funcSpecs []funcspec.FuncSpec) string {
	if strings.HasPrefix(value, "\"") {
		return "string"
	}

	if ssacparser.IsLiteral(value) {
		return inferLiteralType(value)
	}

	parts := strings.SplitN(value, ".", 2)
	if len(parts) < 2 {
		typeName, ok := definedVars[value]
		if !ok {
			return ""
		}
		return typeName
	}
	source, field := parts[0], parts[1]

	switch source {
	case "request":
		return resolveOpenAPIFieldType(doc, funcName, field)
	case "currentUser":
		return ""
	}

	typeName, ok := definedVars[source]
	if !ok {
		return ""
	}

	tableName := modelToTable(typeName)
	if goType := resolveDDLColumnType(st, tableName, field); goType != "" {
		return goType
	}

	return resolveFuncResponseFieldType(funcSpecs, typeName, field)
}
