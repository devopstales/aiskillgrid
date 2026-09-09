package service

import (
	"testing"
)

// TestOpenEmptyProjectIDAborts guards the handle seam: an empty project id
// must abort open with an error and never yield a partial handle.
func TestOpenEmptyProjectIDAborts(t *testing.T) {
	dataDir := t.TempDir()
	svc := New(dataDir)

	for _, id := range []string{"", "   "} {
		h, cleanup, err := svc.Open(id)
		if err == nil {
			t.Fatalf("Open(%q): expected abort, got nil error", id)
		}
		if h != nil {
			t.Fatalf("Open(%q): expected nil handle on abort, got non-nil", id)
		}
		if cleanup != nil {
			cleanup()
		}
	}
}

// TestOpenInvalidProjectIDAborts covers ids store.Open rejects ("..").
func TestOpenInvalidProjectIDAborts(t *testing.T) {
	dataDir := t.TempDir()
	svc := New(dataDir)

	h, cleanup, err := svc.Open("a/../b")
	if err == nil {
		t.Fatal("Open(a/../b): expected abort, got nil error")
	}
	if h != nil {
		t.Fatal("Open(a/../b): expected nil handle on abort, got non-nil")
	}
	if cleanup != nil {
		cleanup()
	}
}
