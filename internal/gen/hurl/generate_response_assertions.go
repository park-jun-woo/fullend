//ff:func feature=gen-hurl type=generator
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

		if prop.Type.Slice()[0] == "array" {
			asserts = append(asserts, fmt.Sprintf("jsonpath %q isCollection", prefix))
		} else if prop.Type.Slice()[0] == "object" {
			asserts = append(asserts, fmt.Sprintf("jsonpath %q exists", prefix))
			// Check nested ID field (use snake_case for sqlc JSON tags).
			if nested := prop.Properties; nested != nil {
				if _, hasID := nested["ID"]; hasID {
					asserts = append(asserts, fmt.Sprintf("jsonpath %q exists", prefix+".id"))
				}
			}
		} else {
			asserts = append(asserts, fmt.Sprintf("jsonpath %q exists", prefix))
		}
	}

	return sortStringSlice(asserts)
}
