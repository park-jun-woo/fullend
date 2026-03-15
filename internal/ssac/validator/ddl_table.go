//ff:type feature=symbol type=model
//ff:what DDL에서 파싱한 테이블 컬럼 정보
package validator

// DDLTable은 DDL에서 파싱한 테이블 컬럼 정보다.
type DDLTable struct {
	Columns     map[string]string // snake_case 컬럼명 → Go 타입
	ColumnOrder []string          // DDL 정의 순서 보존
	ForeignKeys []ForeignKey      // FK 관계 목록
	Indexes     []Index           // 인덱스 목록
	PrimaryKey  []string          // PK 컬럼명 목록 (e.g. ["id"])
}
