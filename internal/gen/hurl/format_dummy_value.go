//ff:func feature=gen-hurl type=util control=selection
//ff:what 더미 값을 JSON 리터럴로 포맷한다
package hurl

import "fmt"

// formatDummyValue formats a dummy value as a JSON literal string.
func formatDummyValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case int:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case map[string]interface{}:
		return "{}"
	case []interface{}:
		return "[]"
	default:
		return fmt.Sprintf("%q", fmt.Sprint(val))
	}
}
