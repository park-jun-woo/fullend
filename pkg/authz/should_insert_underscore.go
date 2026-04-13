//ff:func feature=pkg-authz type=util control=sequence
//ff:what shouldInsertUnderscore — snake_case 변환 시 i 위치에 _ 삽입 여부 판단

package authz

import "unicode"

func shouldInsertUnderscore(runes []rune, i int) bool {
	if i == 0 || !unicode.IsUpper(runes[i]) {
		return false
	}
	prev := runes[i-1]
	nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
	return unicode.IsLower(prev) || nextIsLower
}
