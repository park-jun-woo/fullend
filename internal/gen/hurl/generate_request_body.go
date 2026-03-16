//ff:func feature=gen-hurl type=generator control=iteration
//ff:what 스키마에서 JSON 요청 본문을 빌드한다
package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// generateRequestBody builds a JSON request body from an OpenAPI schema.
// Returns the JSON string and a map of field→value for assertion generation.
func generateRequestBody(schema *openapi3.Schema, checkEnums map[string][]string) (string, map[string]interface{}) {
	if schema == nil {
		return "{}", nil
	}

	values := make(map[string]interface{})
	var lines []string
	for name, propRef := range schema.Properties {
		prop := propRef.Value
		val := generateDummyValue(name, prop, checkEnums)
		values[name] = val
		lines = append(lines, fmt.Sprintf("  %s: %s", formatDummyValue(name), formatDummyValue(val)))
	}

	// Sort for deterministic output.
	sortedLines := sortStringSlice(lines)
	return "{\n" + strings.Join(sortedLines, ",\n") + "\n}", values
}
