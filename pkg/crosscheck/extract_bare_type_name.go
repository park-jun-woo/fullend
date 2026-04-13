//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what extractBareTypeName — Go 타입 문자열에서 `[]`, `*`, `pkg.` 등 장식 제거한 basename

package crosscheck

import "strings"

// extractBareTypeName strips "[]" slice prefix, "*" pointer prefix, and any
// package prefix, returning the basename. e.g. "[]*model.Action" → "Action".
// Also strips generic brackets: "pagination.Cursor[T]" → "Cursor".
func extractBareTypeName(t string) string {
	t = strings.TrimSpace(t)
	for strings.HasPrefix(t, "[]") || strings.HasPrefix(t, "*") {
		if strings.HasPrefix(t, "[]") {
			t = t[2:]
		} else {
			t = t[1:]
		}
	}
	if i := strings.IndexByte(t, '['); i > 0 {
		t = t[:i]
	}
	if i := strings.LastIndexByte(t, '.'); i >= 0 {
		t = t[i+1:]
	}
	return t
}
