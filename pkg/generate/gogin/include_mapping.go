//ff:type feature=gen-gogin type=model
//ff:what resolved x-include to DDL FK mapping (forward FK only)

package gogin

// includeMapping represents a resolved x-include → DDL FK mapping (forward FK only).
type includeMapping struct {
	IncludeName string // "instructor" — derived from FK column (strip _id)
	FieldName   string // "Instructor"
	FieldType   string // "*User"
	FKColumn    string // "instructor_id"
	TargetTable string // "users"
	TargetModel string // "User"
}
