//ff:func feature=gen-hurl type=generator control=iteration dimension=1
//ff:what jsonpath 어설션을 생성한다
package hurl

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// generateResponseAssertions builds jsonpath assertions from a response schema.
// Response field names are converted to snake_case to match sqlc-generated JSON tags.
func generateResponseAssertions(schema *openapi3.Schema, sentValues map[string]interface{}) []string {
	if schema == nil {
		return nil
	}

	var asserts []string
	for name, propRef := range schema.Properties {
		prop := propRef.Value
		prefix := "$." + name
		typeName := prop.Type.Slice()[0]
		if typeName == "array" {
			asserts = append(asserts, fmt.Sprintf("jsonpath %q isCollection", prefix))
			continue
		}
		asserts = append(asserts, fmt.Sprintf("jsonpath %q exists", prefix))
		if typeName == "object" && prop.Properties != nil && prop.Properties["ID"] != nil {
			asserts = append(asserts, fmt.Sprintf("jsonpath %q exists", prefix+".id"))
		}
	}

	return sortStringSlice(asserts)
}
