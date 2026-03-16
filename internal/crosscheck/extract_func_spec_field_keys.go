//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what func spec 응답 필드에서 JSON 키 목록 추출
package crosscheck

import "github.com/geul-org/fullend/internal/funcspec"

// extractFuncSpecFieldKeys extracts JSON field keys from a func spec's response fields.
func extractFuncSpecFieldKeys(fs funcspec.FuncSpec) []string {
	var keys []string
	for _, f := range fs.ResponseFields {
		if f.JSONName != "" {
			keys = append(keys, f.JSONName)
		} else {
			keys = append(keys, f.Name)
		}
	}
	return keys
}
