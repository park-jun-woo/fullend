//ff:func feature=gen-gogin type=util control=sequence
//ff:what removes the model name prefix from a sqlc query name

package gogin

import "strings"

// stripModelPrefix removes the model name prefix from a sqlc query name.
// e.g. "CourseFindByID" with modelName "Course" -> "FindByID".
// If no prefix matches, returns the original name (backward compat).
func stripModelPrefix(queryName, modelName string) string {
	if strings.HasPrefix(queryName, modelName) && len(queryName) > len(modelName) {
		return queryName[len(modelName):]
	}
	return queryName
}
