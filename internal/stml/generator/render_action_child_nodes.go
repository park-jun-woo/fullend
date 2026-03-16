//ff:func feature=stml-gen type=generator control=iteration dimension=1
//ff:what Action 컨텍스트에서 ChildNode를 Kind별로 분기하여 렌더링한다
package generator

import "github.com/geul-org/fullend/internal/stml/parser"

// renderActionChildNodes renders ChildNode slice in DOM order for action context.
func renderActionChildNodes(nodes []parser.ChildNode, formName string, indent int) []string {
	var lines []string
	for _, ch := range nodes {
		switch ch.Kind {
		case "bind":
			lines = append(lines, renderFieldJSX(*ch.Bind, formName, indent))
		case "static":
			lines = append(lines, renderStaticActionJSX(*ch.Static, formName, indent))
		}
	}
	return lines
}
