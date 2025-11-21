package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	t.Run("creates manager with valid path", func(t *testing.T) {
		testConfigDir := t.TempDir()
		m, err := NewManager(testConfigDir)

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
	t.Run("creates manager with invalid path", func(t *testing.T) {

		_, err := NewManager("")

		if err == nil {
			t.Fatal("expected error, got none")
		}

		if !errors.Is(err, ErrMissingHomeDir) {
			t.Errorf("expected %v, got %v", ErrMissingHomeDir, err)
		}
	})
}

func TestManager_GetFilePath(t *testing.T) {
	want := filepath.Join(t.TempDir(), "file.json")
	m := Manager{
		filePath: want,
	}

	got := m.GetFilePath()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestManager_EnsureFolder(t *testing.T) {
	t.Run("create folder when it doesn't exists", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigDir := filepath.Join(tempDir, ".teamtime")

		m, err := NewManager(tempDir)
		if err != nil {
			t.Fatalf("error creating a new manager: %v", err)
		}

		err = m.EnsureFolder()
		if err != nil {
			t.Fatalf("ensureFolder failed %v", err)
		}

		info, err := os.Stat(testConfigDir)

		if os.IsNotExist(err) {
			t.Errorf("expected folder to exists, got error: %v", err)
		}

		if !info.IsDir() {
			t.Error("expected path to be a directory")
		}

		expectedPerm := os.FileMode(0755)
		if info.Mode().Perm() != expectedPerm {
			t.Errorf("expected permission 0775, got %04o", info.Mode().Perm())
		}
	})
	t.Run("succeeds when folder already exists", func(t *testing.T) {
		tempDir := t.TempDir()
		testConfigDir := filepath.Join(tempDir, ".teamtime")

		err := os.MkdirAll(testConfigDir, 0775)

		if err != nil {
			t.Fatalf("failed to create test directory: %v", err)
		}

		m, err := NewManager(tempDir)
		if err != nil {
			t.Fatalf("error creating a new manager: %v", err)
		}

		err = m.EnsureFolder()
		if err != nil {
			t.Fatalf("expected no error when folder exists, got: %v", err)
		}
	})
}
