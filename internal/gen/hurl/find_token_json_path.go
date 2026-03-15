//ff:func feature=gen-hurl type=util
//ff:what 응답 스키마에서 토큰 필드 JSON path를 찾는다
package hurl

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// findTokenJSONPath finds the jsonpath to a string token field in the response schema.
// Handles both flat ("$.access_token") and nested ("$.token.access_token") structures.
func findTokenJSONPath(respSchema *openapi3.Schema) string {
	if respSchema == nil || respSchema.Properties == nil {
		return "token"
	}

	// First pass: look for a string-typed token field at the top level.
	for name := range respSchema.Properties {
		lname := strings.ToLower(name)
		prop := respSchema.Properties[name].Value
		if prop == nil {
			continue
		}
		if (strings.Contains(lname, "token") || strings.Contains(lname, "accesstoken")) &&
			len(prop.Type.Slice()) > 0 && prop.Type.Slice()[0] == "string" {
			return name
		}
	}

	// Second pass: if a token field is an object, look inside for a string token field.
	for name := range respSchema.Properties {
		lname := strings.ToLower(name)
		prop := respSchema.Properties[name].Value
		if prop == nil {
			continue
		}
		if !strings.Contains(lname, "token") {
			continue
		}
		if len(prop.Type.Slice()) > 0 && prop.Type.Slice()[0] == "object" && prop.Properties != nil {
			for innerName := range prop.Properties {
				innerLname := strings.ToLower(innerName)
				innerProp := prop.Properties[innerName].Value
				if innerProp == nil {
					continue
				}
				if (strings.Contains(innerLname, "token") || strings.Contains(innerLname, "access")) &&
					len(innerProp.Type.Slice()) > 0 && innerProp.Type.Slice()[0] == "string" {
					return name + "." + innerName
				}
			}
		}
	}

	// Fallback: return first token-like field regardless of type.
	for name := range respSchema.Properties {
		lname := strings.ToLower(name)
		if strings.Contains(lname, "token") || strings.Contains(lname, "accesstoken") {
			return name
		}
	}

	return "token"
}
