//ff:func feature=gen-react type=generator control=iteration dimension=1
//ff:what api.ts 파일을 생성한다 (fetch 래퍼 함수 포함)

package react

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeAPIClient generates api.ts with fetch wrappers using object parameters.
func writeAPIClient(srcDir string, doc *openapi3.T) error {
	var b strings.Builder
	b.WriteString("const BASE = '/api'\n\n")

	if doc == nil || doc.Paths == nil {
		b.WriteString("export const api = {}\n")
		return os.WriteFile(filepath.Join(srcDir, "api.ts"), []byte(b.String()), 0644)
	}

	endpoints := collectEndpoints(doc)

	for _, ep := range endpoints {
		writeEndpointFunc(&b, ep)
	}

	writeApiNamespace(&b, endpoints)

	return os.WriteFile(filepath.Join(srcDir, "api.ts"), []byte(b.String()), 0644)
}
