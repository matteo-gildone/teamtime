package config

import (
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	t.Run("creates manager with valid path", func(t *testing.T) {
		m, err := NewManager()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if m == nil {
			t.Fatal("expected non-nil manager")
		}

		if m.filePath == "" {
			t.Error("expected filePath to be set")
		}

		if !filepath.IsAbs(m.filePath) {
			t.Error("expected absolute path")
		}

		if filepath.Base(m.filePath) != "colleagues.json" {
			t.Errorf("expected filename colleague.json, got %s", filepath.Base(m.filePath))
		}
	})
}
