//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what Extracts the plural table name from a URL path.
package hurl

import "strings"

// inferTableFromPath extracts the plural table name from a URL path.
// e.g. "/courses/{CourseID}" -> "courses", "/lessons/{LessonID}" -> "lessons"
func inferTableFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Find the first non-parameter segment.
	for _, p := range parts {
		if !strings.HasPrefix(p, "{") && !strings.HasPrefix(p, ":") {
			return p
		}
	}
	return ""
}
