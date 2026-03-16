//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=args-inputs
//ff:what map[string]string을 Go struct 리터럴 필드 문자열로 변환
package generator

import (
	"sort"
	"strings"

	"github.com/ettle/strcase"
)

// buildInputFieldsFromMap은 map[string]string을 Go struct 리터럴 필드로 변환한다.
func buildInputFieldsFromMap(inputs map[string]string) string {
	keys := make([]string, 0, len(inputs))
	for k := range inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var fields []string
	for _, k := range keys {
		fields = append(fields, strcase.ToGoPascal(k)+": "+inputValueToCode(inputs[k]))
	}
	return strings.Join(fields, ", ")
}
