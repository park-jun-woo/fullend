//ff:func feature=contract type=util control=iteration dimension=1
//ff:what key=value 쌍 목록을 파싱하여 디렉티브 필드를 설정한다
package contract

import (
	"fmt"
	"strings"
)

// parseDirectiveFields parses key=value pairs and sets Directive fields.
func parseDirectiveFields(d *Directive, pairs []string) error {
	for _, p := range pairs {
		key, val, ok := strings.Cut(p, "=")
		if !ok {
			return fmt.Errorf("invalid key=value pair %q in directive", p)
		}
		switch key {
		case "ssot":
			d.SSOT = val
		case "contract":
			d.Contract = val
		default:
			return fmt.Errorf("unknown directive field %q", key)
		}
	}
	return nil
}
