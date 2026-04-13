//ff:func feature=gen-hurl type=generator control=iteration dimension=1
//ff:what Builds a login JSON body matching the registered email.
package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// generateLoginBodyWithEmail builds a login JSON body matching the registered email.
func generateLoginBodyWithEmail(schema *openapi3.Schema, emailPrefix string, checkEnums map[string][]string) string {
	if schema == nil {
		return "{}"
	}
	var lines []string
	for name, propRef := range schema.Properties {
		prop := propRef.Value
		lower := strings.ToLower(name)
		var val interface{}
		switch {
		case lower == "email" || (prop != nil && prop.Format == "email"):
			val = emailPrefix + "@test.com"
		default:
			val = generateDummyValue(name, prop, checkEnums)
		}
		lines = append(lines, fmt.Sprintf("  %s: %s", formatDummyValue(name), formatDummyValue(val)))
	}
	sortedLines := sortStringSlice(lines)
	return "{\n" + strings.Join(sortedLines, ",\n") + "\n}"
}
