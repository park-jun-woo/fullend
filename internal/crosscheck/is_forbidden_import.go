//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what import가 금지된 I/O 패키지 접두사와 일치하는지 확인
package crosscheck

import "strings"

// isForbiddenImport checks if an import matches any forbidden I/O package prefix.
func isForbiddenImport(imp string) bool {
	for _, forbidden := range forbiddenImportPrefixes {
		if imp == forbidden || strings.HasPrefix(imp, forbidden+"/") {
			return true
		}
	}
	return false
}
