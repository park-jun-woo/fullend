//ff:type feature=symbol type=model
//ff:what 외부 SSOT에서 수집한 심볼 정보 테이블 + Clone
package validator

// SymbolTable은 외부 SSOT에서 수집한 심볼 정보다.
type SymbolTable struct {
	Models     map[string]ModelSymbol     // "User" → {Methods: {"FindByID": ...}}
	Operations map[string]OperationSymbol // "Login" → {RequestFields, PathParams, ...}
	Funcs      map[string]bool            // "calculateRefund" → true
	DDLTables  map[string]DDLTable        // "users" → {Columns: {"id": "int64", ...}}
	DTOs           map[string]bool              // "Token" → true (DDL 테이블 없는 순수 DTO)
	RequestSchemas map[string]RequestSchema     // operationId → requestBody 필드 제약
}

// Clone returns a shallow copy of the SymbolTable with a deep-copied Models map.
// Other maps (Operations, Funcs, DDLTables, DTOs) share the original data
// since only Models is mutated by injectFuncErrStatus in the gen path.
func (st *SymbolTable) Clone() *SymbolTable {
	clone := &SymbolTable{
		Operations:     st.Operations,
		Funcs:          st.Funcs,
		DDLTables:      st.DDLTables,
		DTOs:           st.DTOs,
		RequestSchemas: st.RequestSchemas,
		Models:     make(map[string]ModelSymbol, len(st.Models)),
	}
	for k, ms := range st.Models {
		methods := make(map[string]MethodInfo, len(ms.Methods))
		for mk, mi := range ms.Methods {
			methods[mk] = mi
		}
		clone.Models[k] = ModelSymbol{Methods: methods}
	}
	return clone
}
