//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what PascalCase 컴포넌트명을 kebab-case 파일명으로 변환한다

package react

// componentToFileName converts "ListCoursesPage" -> "list-courses-page"
func componentToFileName(component string) string {
	var result []byte
	for i, r := range component {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '-')
		}
		result = append(result, byte(r|0x20)) // toLower
	}
	return string(result)
}
