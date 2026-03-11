package authz

import (
	"os"
	"testing"
)

func TestCheckDisabled(t *testing.T) {
	os.Setenv("DISABLE_AUTHZ", "1")
	defer os.Unsetenv("DISABLE_AUTHZ")

	resp, err := Check(CheckRequest{
		Action:   "read",
		Resource: "gig",
		UserID:   1,
	})
	if err != nil {
		t.Fatalf("expected no error with DISABLE_AUTHZ=1, got: %v", err)
	}
	_ = resp
}

func TestCheckNotInitialized(t *testing.T) {
	os.Unsetenv("DISABLE_AUTHZ")
	globalEval = nil

	_, err := Check(CheckRequest{
		Action:   "read",
		Resource: "gig",
		UserID:   1,
	})
	if err == nil {
		t.Fatal("expected error when not initialized")
	}
	if err.Error() != "authz not initialized" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckRequestFields(t *testing.T) {
	req := CheckRequest{
		Action:     "AcceptProposal",
		Resource:   "gig",
		UserID:     42,
		ResourceID: 99,
	}
	if req.Action != "AcceptProposal" {
		t.Fatal("Action mismatch")
	}
	if req.Resource != "gig" {
		t.Fatal("Resource mismatch")
	}
	if req.UserID != 42 {
		t.Fatal("UserID mismatch")
	}
	if req.ResourceID != 99 {
		t.Fatal("ResourceID mismatch")
	}
}
