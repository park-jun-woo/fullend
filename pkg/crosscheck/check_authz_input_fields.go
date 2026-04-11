//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkAuthzInputFields — @auth input 필드가 authz CheckRequest에 있는지 검증
package crosscheck

func checkAuthzInputFields(funcName string, inputs map[string]string) []CrossError {
	var errs []CrossError
	for field := range inputs {
		if !authzCheckRequestFields[field] {
			errs = append(errs, CrossError{Rule: "X-60", Context: funcName, Level: "ERROR",
				Message: "@auth input " + field + " not in authz CheckRequest fields"})
		}
	}
	return errs
}
