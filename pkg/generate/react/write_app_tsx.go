//ff:func feature=gen-react type=generator control=iteration dimension=1
//ff:what App.tsx 파일을 생성한다 (페이지 라우팅 포함)

package react

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeAppTSX generates App.tsx by scanning actual page files on disk
// and matching them to OpenAPI paths via stml page->operationID mapping.
func writeAppTSX(srcDir string, doc *openapi3.T, stmlPages []string, stmlPageOps map[string]string) error {
	pageFiles := scanPageFiles(filepath.Join(srcDir, "pages"))
	opPaths := buildOpPaths(doc)
	uniqueRoutes := buildAppRoutes(pageFiles, stmlPageOps, opPaths)

	var b strings.Builder
	b.WriteString("import { Routes, Route } from 'react-router-dom'\n")
	for _, r := range uniqueRoutes {
		b.WriteString(fmt.Sprintf("import %s from './pages/%s'\n", r.component, r.fileName))
	}
	b.WriteString("\n")

	b.WriteString("export default function App() {\n")
	b.WriteString("  return (\n")
	b.WriteString("    <Routes>\n")
	for _, r := range uniqueRoutes {
		b.WriteString(fmt.Sprintf("      <Route path=\"%s\" element={<%s />} />\n", r.path, r.component))
	}
	b.WriteString("    </Routes>\n")
	b.WriteString("  )\n")
	b.WriteString("}\n")

	return os.WriteFile(filepath.Join(srcDir, "App.tsx"), []byte(b.String()), 0644)
}
