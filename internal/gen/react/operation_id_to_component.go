//ff:func feature=gen-react type=util control=sequence
//ff:what operationID를 컴포넌트명으로 변환한다

package react

// operationIDToComponent converts "ListCourses" -> "ListCoursesPage"
func operationIDToComponent(opID string) string {
	return opID + "Page"
}
