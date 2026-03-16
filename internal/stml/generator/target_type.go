//ff:type feature=stml-gen type=generator
//ff:what 코드 생성 백엔드를 추상화하는 인터페이스
package generator

import "github.com/geul-org/fullend/internal/stml/parser"

// Target abstracts the code generation backend.
// Implement this interface to support a new framework (e.g. Vue, Svelte).
type Target interface {
	GeneratePage(page parser.PageSpec, specsDir string, opts GenerateOptions) string
	FileExtension() string
	Dependencies(pages []parser.PageSpec) map[string]string
}

// compile-time check
var _ Target = (*ReactTarget)(nil)
