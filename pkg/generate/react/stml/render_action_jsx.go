//ff:func feature=stml-gen type=generator control=sequence
//ff:what ActionBlock의 폼 또는 버튼 JSX를 Fields 유무에 따라 생성한다
package stml

import stmlparser "github.com/park-jun-woo/fullend/pkg/parser/stml"

// renderActionJSX generates JSX for an ActionBlock.
func renderActionJSX(a stmlparser.ActionBlock, indent int) string {
	if len(a.Fields) == 0 {
		return renderActionButton(a, indent)
	}
	return renderActionForm(a, indent)
}
