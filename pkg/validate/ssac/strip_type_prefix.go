//ff:func feature=rule type=util control=sequence
//ff:what stripTypePrefix — タイプ名から []、パッケージ接頭辞を除去して純粋な型名を返す
package ssac

import "strings"

func stripTypePrefix(typeName string) string {
	t := strings.TrimPrefix(typeName, "[]")
	if idx := strings.LastIndexByte(t, '.'); idx >= 0 {
		t = t[idx+1:]
	}
	return t
}
