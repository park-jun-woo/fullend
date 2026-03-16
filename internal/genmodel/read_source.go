//ff:func feature=genmodel type=loader control=sequence
//ff:what URL 또는 파일 경로로부터 OpenAPI 소스 데이터를 읽어온다
package genmodel

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func readSource(source string) ([]byte, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return nil, fmt.Errorf("fetch URL: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("fetch URL: status %d", resp.StatusCode)
		}
		return io.ReadAll(resp.Body)
	}
	return os.ReadFile(source)
}
