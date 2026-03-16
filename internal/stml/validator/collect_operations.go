//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what 단일 경로의 모든 HTTP 메서드 오퍼레이션을 SymbolTable에 수집
package validator

func collectOperations(pathItem map[string]openAPIOperation, schemas map[string]openAPISchema, st *SymbolTable) {
	for method, op := range pathItem {
		if op.OperationID == "" {
			continue
		}
		api := buildAPISymbol(method, op, schemas)
		st.Operations[op.OperationID] = api
	}
}
