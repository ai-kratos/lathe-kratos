package store_test

import (
	"testing"
	"time"

	"github.com/devenjarvis/lathe/internal/store"
)

func TestWriteReadMetadata(t *testing.T) {
	dir := t.TempDir()
	original := &store.Tutorial{
		Slug:    "test-tutorial",
		Title:   "Test Tutorial",
		Topic:   "test tutorial",
		Created: time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC),
		Status:  store.StatusVerified,
		Series:  false,
	}

	if err := store.WriteMetadata(dir, original); err != nil {
		t.Fatalf("WriteMetadata() error = %v", err)
	}

	got, err := store.ReadMetadata(dir)
	if err != nil {
		t.Fatalf("ReadMetadata() error = %v", err)
	}
	if got.Slug != original.Slug {
		t.Errorf("Slug = %q, want %q", got.Slug, original.Slug)
	}
	if got.Status != original.Status {
		t.Errorf("Status = %q, want %q", got.Status, original.Status)
	}
}

func TestReadMetadataNotFound(t *testing.T) {
	_, err := store.ReadMetadata("/nonexistent/path/abc123")
	if err == nil {
		t.Error("ReadMetadata() expected error for missing file, got nil")
	}
}
