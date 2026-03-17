//ff:func feature=stml-gen type=generator control=sequence
//ff:what ComponentRef의 JSX를 생성한다
package generator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

// renderComponentJSX generates JSX for a ComponentRef.
func renderComponentJSX(c parser.ComponentRef, dataVar string, indent int) string {
	ind := indentStr(indent)
	if c.Bind != "" {
		return fmt.Sprintf("%s<%s data={%s.%s} />", ind, c.Name, dataVar, c.Bind)
	}
	return fmt.Sprintf("%s<%s />", ind, c.Name)
}
