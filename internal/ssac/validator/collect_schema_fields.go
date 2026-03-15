//ff:func feature=symbol type=util
//ff:what 인라인 properties와 $ref 모두에서 필드를 수집한다
package validator

import "strings"

// collectSchemaFields는 인라인 properties와 $ref 모두에서 필드를 수집한다.
func collectSchemaFields(schema openAPISchema, schemas map[string]openAPISchema) []string {
	var fields []string

	// 인라인 properties
	for k := range schema.Properties {
		fields = append(fields, k)
	}

	// $ref 해결
	if schema.Ref != "" {
		name := schema.Ref[strings.LastIndex(schema.Ref, "/")+1:]
		if resolved, ok := schemas[name]; ok {
			for k := range resolved.Properties {
				fields = append(fields, k)
			}
		}
	}

	return fields
}
