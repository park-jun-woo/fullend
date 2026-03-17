//ff:func feature=stml-gen type=util control=sequence
//ff:what 페이지의 루트 엘리먼트 태그와 클래스명을 결정한다
package generator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

func findRootElement(page parser.PageSpec) (string, string) {
	if len(page.Children) == 1 && page.Children[0].Kind == "static" {
		se := page.Children[0].Static
		return se.Tag, se.ClassName
	}
	return "div", ""
}
