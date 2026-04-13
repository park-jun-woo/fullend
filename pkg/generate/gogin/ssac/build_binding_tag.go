//ff:func feature=ssac-gen type=util control=selection topic=request-params
//ff:what FieldConstraint를 Gin binding 태그 문자열로 변환
package ssac

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func buildBindingTag(fc rule.FieldConstraint) string {
	var parts []string
	if fc.Required {
		parts = append(parts, "required")
	}
	switch fc.Format {
	case "email":
		parts = append(parts, "email")
	case "uuid":
		parts = append(parts, "uuid")
	case "uri":
		parts = append(parts, "uri")
	}
	if fc.MinLength != nil {
		parts = append(parts, fmt.Sprintf("min=%d", *fc.MinLength))
	}
	if fc.MaxLength != nil {
		parts = append(parts, fmt.Sprintf("max=%d", *fc.MaxLength))
	}
	if fc.Minimum != nil {
		parts = append(parts, fmt.Sprintf("gte=%g", *fc.Minimum))
	}
	if fc.Maximum != nil {
		parts = append(parts, fmt.Sprintf("lte=%g", *fc.Maximum))
	}
	if len(fc.Enum) > 0 {
		parts = append(parts, "oneof="+strings.Join(fc.Enum, " "))
	}
	if len(parts) == 0 {
		return ""
	}
	return `binding:"` + strings.Join(parts, ",") + `"`
}
