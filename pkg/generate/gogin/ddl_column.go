//ff:type feature=gen-gogin type=model
//ff:what column parsed from a CREATE TABLE statement

package gogin

// ddlColumn represents a column parsed from a CREATE TABLE statement.
type ddlColumn struct {
	Name      string // e.g. "instructor_id"
	GoName    string // e.g. "InstructorID"
	GoType    string // e.g. "int64"
	FKTable   string // e.g. "users" — REFERENCES target table (empty if no FK)
	NotNull   bool   // true if column has NOT NULL or is PRIMARY KEY
	Sensitive bool   // true if column has -- @sensitive annotation → json:"-"
}
