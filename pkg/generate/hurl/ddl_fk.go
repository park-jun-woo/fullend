//ff:type feature=gen-hurl type=model
//ff:what DDL FK 관계 정보 구조체
package hurl

// ddlFK holds minimal DDL info needed for FK-based delete ordering.
type ddlFK struct {
	TableName string
	FKTables  []string // tables referenced via FOREIGN KEY
}
