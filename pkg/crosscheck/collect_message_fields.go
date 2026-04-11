//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectMessageFields — StructInfo에서 messageType 일치하는 필드명 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectMessageFields(structs []ssac.StructInfo, messageType string) []string {
	for _, st := range structs {
		if st.Name != messageType {
			continue
		}
		var fields []string
		for _, f := range st.Fields {
			fields = append(fields, f.Name)
		}
		return fields
	}
	return nil
}
