//ff:type feature=rule type=model
//ff:what Ground 의 구조적 신 필드 타입 정의 (generate 소비 대상)
package rule

// ModelInfo — Go 인터페이스 + sqlc 쿼리 + FuncSpec @error 가 결합된 모델 메타.
// 키는 일반 모델명("User") 또는 @call 함수 모델("auth._func").
type ModelInfo struct {
	Name    string
	Methods map[string]MethodInfo
}

// MethodInfo — 모델 메서드 하나의 메타.
type MethodInfo struct {
	Cardinality string   // "one" / "many" / "exec" (sqlc 유래 또는 빈값)
	Params      []string // 매개변수 이름 순서 (sqlc 의 $N 또는 iface 시그니처 유래)
	ErrStatus   int      // FuncSpec @error 코드 (기본 0 = 미지정)
}

// TableInfo — DDL 테이블의 컬럼 메타.
// Columns 는 이름→타입 조회용 맵, ColumnOrder 는 원본 DDL 선언 순서 (map 비결정성 해소).
type TableInfo struct {
	Name        string
	Columns     map[string]string
	ColumnOrder []string
}

// OperationInfo — OpenAPI operation 하나의 메타.
// Pagination/Sort/Filter 는 해당 x- 확장이 있을 때만 nil 이 아님.
type OperationInfo struct {
	ID             string
	Method         string
	Path           string
	PathParams     []PathParam
	HasRequestBody bool
	Pagination     *PaginationSpec
	Sort           *SortSpec
	Filter         *FilterSpec
}

// PathParam — operation 의 경로 파라미터.
type PathParam struct {
	Name   string
	GoType string
}

// PaginationSpec — x-pagination.
type PaginationSpec struct {
	Style        string // "offset" / "cursor"
	DefaultLimit int
	MaxLimit     int
}

// SortSpec — x-sort.
type SortSpec struct {
	Allowed   []string
	Default   string
	Direction string
}

// FilterSpec — x-filter.
type FilterSpec struct {
	Allowed []string
}

// RequestSchemaInfo — OpenAPI requestBody 의 필드 제약 집합.
type RequestSchemaInfo struct {
	Fields map[string]FieldConstraint
}

// FieldConstraint — 필드 하나의 JSON Schema 제약.
type FieldConstraint struct {
	Required  bool
	Format    string
	MinLength *int
	MaxLength *int
	Minimum   *float64
	Maximum   *float64
	Pattern   string
	Enum      []string
}
