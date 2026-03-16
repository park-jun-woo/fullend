//ff:func feature=gen-hurl type=generator control=iteration dimension=1
//ff:what Builds a JSON body with role and email overrides, resolving FK fields to captured variables.
package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// generateRequestBodyWithOverrides builds a JSON body with role and email overrides.
// FK fields (ending in _id) are resolved to captured variable references if available.
func generateRequestBodyWithOverrides(schema *openapi3.Schema, role, emailPrefix string, checkEnums map[string][]string, captures map[string]bool) string {
	if schema == nil {
		return "{}"
	}
	var lines []string
	for name, propRef := range schema.Properties {
		prop := propRef.Value
		lower := strings.ToLower(name)
		capVar := ""
		if strings.HasSuffix(lower, "_id") {
			capVar = findMatchingCapture(name, captures)
		}
		switch {
		case lower == "role" && role != "":
			lines = append(lines, fmt.Sprintf("  %s: %s", formatDummyValue(name), formatDummyValue(role)))
		case lower == "email" || (prop != nil && prop.Format == "email"):
			lines = append(lines, fmt.Sprintf("  %s: %s", formatDummyValue(name), formatDummyValue(emailPrefix+"@test.com")))
		case capVar != "":
			lines = append(lines, fmt.Sprintf("  %s: {{%s}}", formatDummyValue(name), capVar))
		default:
			val := generateDummyValue(name, prop, checkEnums)
			lines = append(lines, fmt.Sprintf("  %s: %s", formatDummyValue(name), formatDummyValue(val)))
		}
	}
	sortedLines := sortStringSlice(lines)
	return "{\n" + strings.Join(sortedLines, ",\n") + "\n}"
}
