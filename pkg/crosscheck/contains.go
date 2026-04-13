//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-ddl
//ff:what contains — string slice 에 대상 포함 여부

package crosscheck

func contains(xs []string, target string) bool {
	for _, x := range xs {
		if x == target {
			return true
		}
	}
	return false
}
