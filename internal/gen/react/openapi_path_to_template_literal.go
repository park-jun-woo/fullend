//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what OpenAPI 경로를 JS 템플릿 리터럴로 변환한다

package react

import "strings"

// openAPIPathToTemplateLiteral converts "/courses/{CourseID}" -> "/courses/${courseID}"
func openAPIPathToTemplateLiteral(path string) string {
	var b strings.Builder
	i := 0
	for i < len(path) {
		start := strings.Index(path[i:], "{")
		if start < 0 {
			b.WriteString(path[i:])
			break
		}
		b.WriteString(path[i : i+start])
		end := strings.Index(path[i+start:], "}")
		if end < 0 {
			b.WriteString(path[i+start:])
			break
		}
		paramName := path[i+start+1 : i+start+end]
		b.WriteString("${" + lcFirst(paramName) + "}")
		i = i + start + end + 1
	}
	return b.String()
}
