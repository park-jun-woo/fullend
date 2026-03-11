package config

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	os.Setenv("TEST_CONFIG_KEY", "hello")
	defer os.Unsetenv("TEST_CONFIG_KEY")

	if v := Get("TEST_CONFIG_KEY"); v != "hello" {
		t.Fatalf("expected 'hello', got %q", v)
	}
}

func TestGetEmpty(t *testing.T) {
	os.Unsetenv("TEST_CONFIG_MISSING")
	if v := Get("TEST_CONFIG_MISSING"); v != "" {
		t.Fatalf("expected empty string, got %q", v)
	}
}

func TestMustGet(t *testing.T) {
	os.Setenv("TEST_CONFIG_MUST", "value")
	defer os.Unsetenv("TEST_CONFIG_MUST")

	if v := MustGet("TEST_CONFIG_MUST"); v != "value" {
		t.Fatalf("expected 'value', got %q", v)
	}
}

func TestMustGetPanics(t *testing.T) {
	os.Unsetenv("TEST_CONFIG_PANIC")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustGet("TEST_CONFIG_PANIC")
}
