//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=args-inputs
//ff:what Inputs map의 value를 paramOrder 또는 알파벳순으로 positional 함수 인자로 변환
package generator

import "strings"

// buildArgsCodeFromInputs는 Inputs map의 value만 추출하여 positional 함수 인자로 변환한다.
// paramOrder가 있으면 그 순서로 배치하고, 없으면 알파벳순 fallback.
func buildArgsCodeFromInputs(inputs map[string]string, paramOrder []string) string {
	if len(inputs) == 0 {
		return ""
	}

	keys := orderInputKeys(inputs, paramOrder)

	var parts []string
	for _, k := range keys {
		parts = append(parts, inputValueToCode(inputs[k]))
	}
	return strings.Join(parts, ", ")
}
