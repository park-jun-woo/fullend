//ff:type feature=gen-gogin type=model
//ff:what parsed CREATE TABLE definition

package gogin

// ddlTable represents a parsed CREATE TABLE definition.
type ddlTable struct {
	TableName string      // e.g. "courses"
	ModelName string      // e.g. "Course"
	Columns   []ddlColumn // ordered columns
}
