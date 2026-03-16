//ff:func feature=symbol type=loader control=sequence
//ff:what 프로젝트 디렉토리에서 심볼 테이블을 구성한다
package validator

import (
	"fmt"
	"path/filepath"
)

// LoadSymbolTable은 프로젝트 디렉토리에서 심볼 테이블을 구성한다.
// 디렉토리 구조:
//
//	<root>/db/queries/*.sql  — sqlc 쿼리 (모델+메서드)
//	<root>/api/openapi.yaml  — OpenAPI spec (request/response)
//	<root>/model/*.go        — Go interface, func
func LoadSymbolTable(root string) (*SymbolTable, error) {
	st := &SymbolTable{
		Models:     make(map[string]ModelSymbol),
		Operations: make(map[string]OperationSymbol),
		Funcs:      make(map[string]bool),
		DDLTables:  make(map[string]DDLTable),
		DTOs:           make(map[string]bool),
		RequestSchemas: make(map[string]RequestSchema),
	}

	if err := st.loadDDL(filepath.Join(root, "db")); err != nil {
		return nil, fmt.Errorf("DDL 로드 실패: %w", err)
	}
	if err := st.loadSqlcQueries(filepath.Join(root, "db", "queries")); err != nil {
		return nil, fmt.Errorf("sqlc 쿼리 로드 실패: %w", err)
	}
	if err := st.loadOpenAPI(filepath.Join(root, "api", "openapi.yaml")); err != nil {
		return nil, fmt.Errorf("OpenAPI 로드 실패: %w", err)
	}
	if err := st.loadGoInterfaces(filepath.Join(root, "model")); err != nil {
		return nil, fmt.Errorf("Go interface 로드 실패: %w", err)
	}

	return st, nil
}
