//ff:func feature=ssac-gen type=util control=selection
//ff:what dotted target의 필드 타입을 DDL/func 출처별로 조회하여 반환
package generator

import "strings"

// ResolveFieldType은 dotted target의 필드 타입을 반환한다.
// "cr.Balance" -> "int64", "wf.OrgID" -> "int64", "wf" -> "" (변수 자체는 default nil 유지)
func (r *FieldTypeResolver) ResolveFieldType(target string) string {
	parts := strings.SplitN(target, ".", 2)
	if len(parts) < 2 {
		return ""
	}
	varName, fieldName := parts[0], parts[1]
	src, ok := r.vars[varName]
	if !ok {
		return ""
	}
	switch src.Kind {
	case "ddl":
		return r.resolveFromDDL(src.ModelName, fieldName)
	case "func":
		return r.resolveFromFunc(src.ModelName, fieldName)
	}
	return ""
}
