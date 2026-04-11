//ff:type feature=manifest type=model
//ff:what Table — DDL CREATE TABLE에서 추출한 테이블 메타데이터
package ddl

// Table holds parsed metadata for a single DDL table.
type Table struct {
	Name        string
	Columns     map[string]string   // column_name → Go type
	ColumnOrder []string            // DDL definition order
	ForeignKeys []ForeignKey
	Indexes     []Index
	PrimaryKey  []string
	VarcharLen  map[string]int      // column → VARCHAR(N)
	CheckEnums  map[string][]string // column → CHECK IN values
}
