//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=args-inputs
//ff:what auth Inputs에서 UserID, Role을 제외한 필드만 필터링
package generator

func filterAuthInputs(inputs map[string]string) map[string]string {
	filtered := make(map[string]string)
	for k, v := range inputs {
		if k != "UserID" && k != "Role" {
			filtered[k] = v
		}
	}
	return filtered
}
