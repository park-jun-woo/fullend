//ff:func feature=gen-gogin type=test control=iteration dimension=1 topic=query-opts
//ff:what generateQueryOpts: 커서 기반 WHERE절 생성 코드 포함 검증

package gogin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateQueryOpts_CursorWhereClause(t *testing.T) {
	dir := t.TempDir()
	if err := generateQueryOpts(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "queryopts.go"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(data)

	// Verify cursor WHERE clause generation.
	checks := []struct {
		name    string
		snippet string
	}{
		{"cursor WHERE block", `if opts.Cursor != ""`},
		{"cursor less-than operator", `op := "<"`},
		{"cursor asc greater-than", `op = ">"`},
		{"cursor column from SortCol", `cursorCol := opts.SortCol`},
		{"cursor default id", `cursorCol = "id"`},
		{"offset skip in cursor mode", `opts.Offset > 0 && opts.Cursor == ""`},
		{"cursor sort fixed comment", "Cursor mode: fixed sort"},
		{"cursor default id DESC", `opts.SortCol = "id"`},
	}

	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if !strings.Contains(src, c.snippet) {
				t.Errorf("generated queryopts.go missing %q", c.snippet)
			}
		})
	}
}
