//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=import-collect
//ff:what 생성된 코드 본문에서 실제 참조되는 패키지만 필터링
package ssac

import "strings"

// filterUsedImports는 생성된 코드 본문에서 실제 참조되는 패키지만 남긴다.
func filterUsedImports(imports []string, body string) []string {
	var used []string
	for _, imp := range imports {
		pkg := imp
		if idx := strings.LastIndex(imp, "/"); idx >= 0 {
			pkg = imp[idx+1:]
		}
		if strings.Contains(body, pkg+".") {
			used = append(used, imp)
		}
	}
	return used
}
