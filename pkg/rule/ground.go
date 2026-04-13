//ff:type feature=rule type=model
//ff:what Ground — 검증 규칙 및 코드 생성이 공유하는 조회 컨텍스트
package rule

// Ground holds lookup data for validation rules and code generation.
// Populated by the caller, stored in ctx via ctx.Set("ground", g).
type Ground struct {
	// 기존 — 주로 validate/crosscheck 소비
	Lookup  map[string]StringSet // "target.kind" -> set of names
	Types   map[string]string    // "target.kind.name" -> type string
	Pairs   map[string]StringSet // "target.pairKind" -> set of "key:value"
	Config  map[string]bool      // config key -> present
	Vars    StringSet            // declared variable names
	Flags   StringSet            // flags for defeaters
	Schemas map[string][]string  // "target.schema" -> ordered field list

	// 신규 — 주로 generate 소비 (Phase002)
	Models     map[string]ModelInfo         // "User" → {Methods}
	Tables     map[string]TableInfo         // "users" → {Columns, ColumnOrder}
	Ops        map[string]OperationInfo     // operationID → operation 메타
	ReqSchemas map[string]RequestSchemaInfo // operationID → requestBody 필드 제약
}
