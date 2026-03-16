//ff:func feature=ssac-gen type=generator control=iteration dimension=1
//ff:what publish의 Options를 Go 코드(WithDelay, WithPriority)로 변환
package generator

import (
	"sort"
	"strings"
)

// buildPublishOptions는 publish의 Options를 Go 코드로 변환한다.
func buildPublishOptions(options map[string]string) string {
	if len(options) == 0 {
		return ""
	}
	keys := make([]string, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, publishOptionToCode(k, options[k]))
	}
	parts = filterNonEmpty(parts)
	if len(parts) == 0 {
		return ""
	}
	return ", " + strings.Join(parts, ", ")
}
