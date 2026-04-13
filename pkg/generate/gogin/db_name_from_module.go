//ff:func feature=gen-gogin type=util control=sequence topic=output
//ff:what dbNameFromModule — modulePath 의 마지막 세그먼트를 기본 DB 이름으로 사용

package gogin

import (
	"path"
	"strings"
)

// dbNameFromModule derives a default database name from a Go module path.
// e.g. "github.com/example/gigbridge" → "gigbridge".
// "example/zenflow/backend" → "backend" 이 되는 것을 피해 마지막 비-"backend" 세그먼트 선택.
func dbNameFromModule(modulePath string) string {
	clean := strings.TrimSuffix(modulePath, "/backend")
	base := path.Base(clean)
	if base == "." || base == "/" || base == "" {
		return "app"
	}
	return base
}
