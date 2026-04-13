//ff:func feature=contract type=util control=sequence
//ff:what 입력 문자열의 SHA-256 해시 앞 7자를 반환한다
package contract

import (
	"crypto/sha256"
	"fmt"
)

// Hash7 computes a 7-character SHA-256 hash.
func Hash7(input string) string {
	h := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", h)[:7]
}
