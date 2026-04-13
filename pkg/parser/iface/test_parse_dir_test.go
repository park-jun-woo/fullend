//ff:func feature=iface-parse type=parser control=sequence
//ff:what ParseDir 기본 동작 검증 — 인터페이스 추출 + 메서드 순서 보존
package iface

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDirExtractsInterfaces(t *testing.T) {
	dir := t.TempDir()
	src := `package model

type UserModel interface {
	Create(email string) error
	FindByID(id int64) (*User, error)
	List(limit int) ([]User, error)
}

type CourseModel interface {
	Get(id int64) (*Course, error)
}
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	ifaces, diags := ParseDir(dir)
	if len(diags) > 0 {
		t.Fatalf("unexpected diagnostics: %v", diags[0].Message)
	}
	if len(ifaces) != 2 {
		t.Fatalf("expected 2 interfaces, got %d", len(ifaces))
	}

	got := map[string][]string{}
	for _, i := range ifaces {
		got[i.Name] = i.Methods
	}

	user, ok := got["UserModel"]
	if !ok {
		t.Fatal("UserModel not found")
	}
	want := []string{"Create", "FindByID", "List"}
	if len(user) != len(want) {
		t.Fatalf("UserModel methods: got %v, want %v", user, want)
	}
	for i, m := range want {
		if user[i] != m {
			t.Errorf("UserModel[%d]: got %q, want %q", i, user[i], m)
		}
	}

	course, ok := got["CourseModel"]
	if !ok || len(course) != 1 || course[0] != "Get" {
		t.Errorf("CourseModel: got %v, want [Get]", course)
	}
}

func TestParseDirMissingDirectory(t *testing.T) {
	ifaces, diags := ParseDir(filepath.Join(t.TempDir(), "nonexistent"))
	if len(ifaces) != 0 || len(diags) != 0 {
		t.Errorf("expected empty + no diags, got %d ifaces / %d diags", len(ifaces), len(diags))
	}
}

func TestParseDirSkipsNonInterface(t *testing.T) {
	dir := t.TempDir()
	src := `package model

type User struct {
	ID int64
}

type ServiceImpl struct{}

type EmptyIface interface{}

func standaloneFunc() {}
`
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	ifaces, diags := ParseDir(dir)
	if len(diags) > 0 {
		t.Fatalf("unexpected diagnostics: %v", diags[0].Message)
	}
	// EmptyIface has no methods → skipped. struct/func 는 인터페이스 아니라 제외.
	if len(ifaces) != 0 {
		t.Errorf("expected 0 interfaces, got %d: %v", len(ifaces), ifaces)
	}
}

func TestParseFileInvalidSource(t *testing.T) {
	dir := t.TempDir()
	src := `not valid go source`
	path := filepath.Join(dir, "broken.go")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	_, diags := ParseDir(dir)
	if len(diags) == 0 {
		t.Fatal("expected diagnostic for invalid source")
	}
}
