//ff:func feature=ssac-gen type=util control=selection topic=publish
//ff:what publish 옵션 키를 Go 코드로 변환 (delay, priority)
package ssac

import "fmt"

func publishOptionToCode(key, value string) string {
	switch key {
	case "delay":
		return fmt.Sprintf("queue.WithDelay(%s)", value)
	case "priority":
		return fmt.Sprintf("queue.WithPriority(%q)", value)
	default:
		return ""
	}
}
