//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-ddl
//ff:what seedHasID — seed 행 중 지정 컬럼 값이 일치하는 것이 있는지

package crosscheck

func seedHasID(seeds []map[string]string, idCol, wantVal string) bool {
	for _, row := range seeds {
		if row[idCol] == wantVal {
			return true
		}
	}
	return false
}
