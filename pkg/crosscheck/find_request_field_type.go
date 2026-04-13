//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what findRequestFieldType — funcspec RequestFields 에서 주어진 field 의 Go 타입 반환

package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/funcspec"

// findRequestFieldType looks up a field's Go type in a FuncSpec's RequestFields.
func findRequestFieldType(fn *funcspec.FuncSpec, field string) string {
	if fn == nil {
		return ""
	}
	for _, f := range fn.RequestFields {
		if f.Name == field {
			return f.Type
		}
	}
	return ""
}
