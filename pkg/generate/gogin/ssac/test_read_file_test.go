//ff:func feature=ssac-gen type=test-helper control=sequence
//ff:what 파일을 읽어 문자열로 반환하는 테스트 헬퍼
package ssac

import (
	"os"
	"testing"
)

func readFile(t *testing.T, path string) (string, error) {
	t.Helper()
	data, err := os.ReadFile(path)
	return string(data), err
}
