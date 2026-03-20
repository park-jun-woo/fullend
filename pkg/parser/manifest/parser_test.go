package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	content := `apiVersion: fullend/v1
kind: Project
metadata:
  name: testapp
backend:
  module: github.com/test/testapp
  auth:
    type: jwt
    claims:
      UserID: "user_id:int64"
      Role: "role"
frontend:
  framework: react
`
	os.WriteFile(filepath.Join(dir, "fullend.yaml"), []byte(content), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.APIVersion != "fullend/v1" {
		t.Errorf("APIVersion = %q, want %q", cfg.APIVersion, "fullend/v1")
	}
	if cfg.Metadata.Name != "testapp" {
		t.Errorf("Metadata.Name = %q, want %q", cfg.Metadata.Name, "testapp")
	}
	if cfg.Backend.Module != "github.com/test/testapp" {
		t.Errorf("Backend.Module = %q", cfg.Backend.Module)
	}
	if cfg.Backend.Auth == nil {
		t.Fatal("Backend.Auth is nil")
	}
	if cfg.Backend.Auth.Type != "jwt" {
		t.Errorf("Auth.Type = %q, want %q", cfg.Backend.Auth.Type, "jwt")
	}
	if len(cfg.Backend.Auth.Claims) != 2 {
		t.Errorf("Claims count = %d, want 2", len(cfg.Backend.Auth.Claims))
	}
	if def, ok := cfg.Backend.Auth.Claims["UserID"]; !ok || def.Key != "user_id" || def.GoType != "int64" {
		t.Errorf("Claims[UserID] = %+v", def)
	}
	if def, ok := cfg.Backend.Auth.Claims["Role"]; !ok || def.Key != "role" || def.GoType != "string" {
		t.Errorf("Claims[Role] = %+v", def)
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := Load("/nonexistent/dir")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "fullend.yaml"), []byte(":\ninvalid: [yaml"), 0644)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_Minimal(t *testing.T) {
	dir := t.TempDir()
	content := `apiVersion: fullend/v1
kind: Project
metadata:
  name: minimal
backend:
  module: github.com/test/minimal
`
	os.WriteFile(filepath.Join(dir, "fullend.yaml"), []byte(content), 0644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Backend.Auth != nil {
		t.Errorf("expected nil Auth, got %+v", cfg.Backend.Auth)
	}
}
