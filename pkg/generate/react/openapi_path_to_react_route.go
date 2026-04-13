//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what OpenAPI 경로를 React Router 경로로 변환한다

package react

import "strings"

// openAPIPathToReactRoute converts "/courses/{CourseID}" -> "/courses/:courseID"
func openAPIPathToReactRoute(path string) string {
	result := path
	for {
		start := strings.Index(result, "{")
		if start < 0 {
			break
		}
		end := strings.Index(result, "}")
		if end < 0 {
			break
		}
		paramName := result[start+1 : end]
		result = result[:start] + ":" + lcFirst(paramName) + result[end+1:]
	}
	return result
}
