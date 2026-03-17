//ff:func feature=stml-gen type=generator control=sequence
//ff:what data-bind 필드의 JSX를 생성한다
package generator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

// renderBindJSX generates JSX for a data-bind field.
func renderBindJSX(b parser.FieldBind, dataVar string, indent int) string {
	ind := indentStr(indent)
	tag := orDefault(b.Tag, "span")
	cls := clsAttr(b.ClassName)
	return fmt.Sprintf("%s<%s%s>{%s.%s}</%s>", ind, tag, cls, dataVar, b.Name, tag)
}
