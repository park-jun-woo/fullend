//ff:func feature=symbol type=util control=sequence topic=openapi
//ff:what $ref가 있으면 해석하여 실제 스키마를 반환
package validator

import "strings"

func resolveSchema(schema openAPISchema, allSchemas map[string]openAPISchema) openAPISchema {
	if schema.Ref == "" {
		return schema
	}
	name := schema.Ref[strings.LastIndex(schema.Ref, "/")+1:]
	if resolved, ok := allSchemas[name]; ok {
		return resolved
	}
	return schema
}
