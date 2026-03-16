//ff:func feature=ssac-parse type=parser control=iteration dimension=1
//ff:what 주석 리스트에서 v2 시퀀스를 추출
package parser

import (
	"go/ast"
	"strings"
)

// parseComments는 주석 리스트에서 v2 시퀀스를 추출한다.
func parseComments(comments []*ast.Comment) ([]Sequence, error) {
	cp := &commentParser{}
	for _, c := range comments {
		line := strings.TrimPrefix(c.Text, "//")
		line = strings.TrimSpace(line)
		if err := cp.processLine(line); err != nil {
			return nil, err
		}
	}
	return cp.sequences, nil
}
