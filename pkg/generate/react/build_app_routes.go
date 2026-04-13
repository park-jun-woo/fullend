//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what 페이지 파일 목록과 OpenAPI 매핑으로 라우트 목록을 생성한다

package react

// buildAppRoutes builds unique sorted routes from page files, stml page->operationID mapping,
// and OpenAPI operationID->path mapping.
func buildAppRoutes(pageFiles []string, stmlPageOps map[string]string, opPaths map[string]string) []route {
	var routes []route
	for _, fileName := range pageFiles {
		matchedPath := resolveRoutePath(fileName, stmlPageOps, opPaths)
		routes = append(routes, route{
			path:      matchedPath,
			component: fileNameToComponent(fileName),
			fileName:  fileName,
		})
	}
	return deduplicateRoutes(routes)
}
