//ff:func feature=gen-react type=util control=sequence
//ff:what 파일명을 PascalCase 컴포넌트명으로 변환한다

package react

import "github.com/ettle/strcase"

// fileNameToComponent converts "course-list-page" -> "CourseListPage"
func fileNameToComponent(fileName string) string {
	return strcase.ToGoPascal(fileName)
}
