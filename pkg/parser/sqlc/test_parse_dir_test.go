//ff:func feature=sqlc-parse type=parser control=sequence
//ff:what ParseDir 기본 동작 검증 — 쿼리 이름·cardinality·params 추출
package sqlc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDirExtractsQueries(t *testing.T) {
	dir := t.TempDir()
	src := `-- name: UserFindByID :one
SELECT * FROM users WHERE id = $1;

-- name: UserList :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UserCreate :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING *;
`
	if err := os.WriteFile(filepath.Join(dir, "users.sql"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	queries, diags := ParseDir(dir)
	if len(diags) > 0 {
		t.Fatalf("unexpected diagnostics: %v", diags[0].Message)
	}
	if len(queries) != 3 {
		t.Fatalf("expected 3 queries, got %d", len(queries))
	}

	// 모델명 동일
	for _, q := range queries {
		if q.Model != "User" {
			t.Errorf("Model: got %q, want %q", q.Model, "User")
		}
	}

	// FindByID — 모델 prefix 제거 + cardinality :one + param 1개
	// strcase.ToGoPascal 은 "id" → "ID" (Go initialism)
	assertQuery(t, queries[0], "FindByID", "one", []string{"ID"})

	// List — params 없음
	assertQuery(t, queries[1], "List", "many", nil)

	// Create — INSERT 컬럼 순서
	assertQuery(t, queries[2], "Create", "one", []string{"Email", "PasswordHash"})
}

func TestParseDirMissingDirectory(t *testing.T) {
	queries, diags := ParseDir(filepath.Join(t.TempDir(), "nonexistent"))
	if len(queries) != 0 || len(diags) != 0 {
		t.Errorf("expected empty + no diags, got %d queries / %d diags", len(queries), len(diags))
	}
}

func TestSqlFileToModel(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"users.sql", "User"},
		{"reservations.sql", "Reservation"},
		{"course_enrollments.sql", "CourseEnrollment"},
		{"single.sql", "Single"},
	}
	for _, tt := range tests {
		got := sqlFileToModel(tt.filename)
		if got != tt.want {
			t.Errorf("sqlFileToModel(%q) = %q, want %q", tt.filename, got, tt.want)
		}
	}
}

func TestStripModelPrefix(t *testing.T) {
	tests := []struct {
		query, model, want string
	}{
		{"UserFindByID", "User", "FindByID"},
		{"FindByID", "User", "FindByID"},
		{"Userdata", "User", "Userdata"}, // prefix 는 같지만 뒷글자가 소문자
		{"UserFindByIDUser", "User", "FindByIDUser"},
	}
	for _, tt := range tests {
		got := stripModelPrefix(tt.query, tt.model)
		if got != tt.want {
			t.Errorf("stripModelPrefix(%q, %q) = %q, want %q", tt.query, tt.model, got, tt.want)
		}
	}
}

func assertQuery(t *testing.T, q Query, name, card string, params []string) {
	t.Helper()
	if q.Name != name {
		t.Errorf("Name: got %q, want %q", q.Name, name)
	}
	if q.Cardinality != card {
		t.Errorf("Cardinality: got %q, want %q", q.Cardinality, card)
	}
	if len(q.Params) != len(params) {
		t.Errorf("Params: got %v, want %v", q.Params, params)
		return
	}
	for i, p := range params {
		if q.Params[i] != p {
			t.Errorf("Params[%d]: got %q, want %q", i, q.Params[i], p)
		}
	}
}
