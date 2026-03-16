//ff:func feature=crosscheck type=rule control=sequence topic=openapi-ddl
//ff:what OpenAPI requestBody 검증 제약 누락·불일치를 DDL 기준으로 검출
package crosscheck

func CheckOpenAPIConstraints(input *CrossValidateInput) []CrossError {
	var errs []CrossError

	st := input.SymbolTable
	if st == nil || st.RequestSchemas == nil {
		return errs
	}

	errs = append(errs, checkC1RequiredMissing(input.ServiceFuncs, st)...)
	errs = append(errs, checkFieldConstraints(st)...)

	return errs
}
