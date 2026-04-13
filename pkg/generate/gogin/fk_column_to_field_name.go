//ff:func feature=gen-gogin type=util control=sequence
//ff:what converts a FK column name to a Go struct field name

package gogin

import "strings"

// fkColumnToFieldName converts a FK column name to a Go struct field name.
// "instructor_id" -> "Instructor", "course_id" -> "Course"
func fkColumnToFieldName(colName string) string {
	name := colName
	if strings.HasSuffix(name, "_id") {
		name = name[:len(name)-3]
	}
	return snakeToGo(name)
}
