//ff:func feature=stml-gen type=generator control=iteration dimension=1
//ff:what FetchBlock의 flat 슬라이스(Binds, Eaches, States, Components) JSX를 생성한다
package generator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

// renderFetchJSXFlatChildren renders flat slices for backward compatibility.
func renderFetchJSXFlatChildren(f parser.FetchBlock, alias string, indent int) []string {
	var lines []string
	for _, b := range f.Binds {
		lines = append(lines, renderBindJSX(b, alias, indent))
	}
	for _, e := range f.Eaches {
		lines = append(lines, renderEachJSX(e, alias, indent))
	}
	for _, s := range f.States {
		lines = append(lines, renderStateJSX(s, alias, indent))
	}
	for _, c := range f.Components {
		lines = append(lines, renderComponentJSX(c, alias, indent))
	}
	return lines
}
