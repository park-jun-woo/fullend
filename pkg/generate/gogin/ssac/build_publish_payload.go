//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=publish
//ff:what publish의 Inputs를 map[string]any 리터럴 필드로 변환
package ssac

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ettle/strcase"
)

// buildPublishPayload는 publish의 Inputs를 map[string]any 리터럴 필드로 변환한다.
func buildPublishPayload(inputs map[string]string) string {
	keys := make([]string, 0, len(inputs))
	for k := range inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var fields []string
	for _, k := range keys {
		fields = append(fields, fmt.Sprintf("\t\t%q: %s,", strcase.ToGoPascal(k), inputValueToCode(inputs[k])))
	}
	return strings.Join(fields, "\n")
}
