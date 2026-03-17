//ff:func feature=ssac-gen type=util control=sequence topic=request-params
//ff:what JSON 필드명과 RequestSchema에서 struct 태그 문자열을 조립
package generator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func buildFieldTag(fieldName string, rs *validator.RequestSchema) string {
	tag := fmt.Sprintf("json:\"%s\"", fieldName)
	if rs == nil {
		return tag
	}
	fc, ok := rs.Fields[fieldName]
	if !ok {
		return tag
	}
	if bt := buildBindingTag(fc); bt != "" {
		tag += " " + bt
	}
	return tag
}
