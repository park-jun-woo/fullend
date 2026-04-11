//ff:func feature=rule type=util control=iteration dimension=1
//ff:what hasStructDef — StructInfo 목록에서 이름이 일치하는 struct 존재 여부
package ssac

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func hasStructDef(structs []parsessac.StructInfo, name string) bool {
	for _, s := range structs {
		if s.Name == name {
			return true
		}
	}
	return false
}
