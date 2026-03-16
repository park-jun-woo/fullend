//ff:func feature=contract type=util control=sequence
//ff:what 모델 구현 메서드의 계약 해시를 계산한다
package contract

import "strings"

// HashModelMethod computes a contract hash for a model implementation method.
// Based on: method name + parameter types + return types.
func HashModelMethod(name string, params []string, returns []string) string {
	parts := []string{name}
	parts = append(parts, strings.Join(params, ","))
	parts = append(parts, strings.Join(returns, ","))
	return Hash7(strings.Join(parts, "|"))
}
