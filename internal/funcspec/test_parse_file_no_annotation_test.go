//ff:func feature=funcspec type=test control=sequence
//ff:what ParseFile: @func 어노테이션 없는 파일에서 nil 반환 검증

package funcspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFileNoAnnotation(t *testing.T) {
	dir := t.TempDir()
	src := `package foo

func Foo() {}
`
	path := filepath.Join(dir, "foo.go")
	os.WriteFile(path, []byte(src), 0644)

	spec, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec != nil {
		t.Error("expected nil for file without @func annotation")
	}
}
