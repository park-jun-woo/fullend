//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkGhostFields — 단일 operation의 response field가 DDL 컬럼에 있는지 검증
package crosscheck

func checkGhostFields(op, table string, fields []string, cols map[string]string) []CrossError {
	var errs []CrossError
	for _, f := range fields {
		if _, ok := cols[f]; ok {
			continue
		}
		if f == "id" {
			continue
		}
		errs = append(errs, CrossError{Rule: "X-9", Context: op, Level: "WARNING",
			Message: "OpenAPI property " + f + " not found in DDL table " + table})
	}
	return errs
}
