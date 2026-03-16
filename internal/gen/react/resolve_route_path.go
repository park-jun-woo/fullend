//ff:func feature=gen-react type=util control=sequence
//ff:what 파일명에서 stml 매핑 또는 파일명 기반으로 라우트 경로를 결정한다

package react

import "strings"

// resolveRoutePath resolves the route path for a page file using stml mapping or filename fallback.
func resolveRoutePath(fileName string, stmlPageOps map[string]string, opPaths map[string]string) string {
	opID := stmlPageOps[fileName]
	if apiPath, ok := opPaths[opID]; ok {
		return openAPIPathToReactRoute(apiPath)
	}
	return "/" + strings.ReplaceAll(strings.TrimSuffix(fileName, "-page"), "-", "/")
}
