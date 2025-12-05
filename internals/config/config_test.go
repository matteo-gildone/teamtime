package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/matteo-gildone/teamtime/internals/types"
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

func TestManager_GetRelativeFilePath(t *testing.T) {

	homeDir := t.TempDir()
	want := "~/file.json"

	m := Manager{
		homeDir:  homeDir,
		filePath: filepath.Join(homeDir, "file.json"),
	}

	got := m.GetRelativeFilePath()

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

		if runtime.GOOS != "windows" {
			expectedPerm := os.FileMode(0755)
			if info.Mode().Perm() != expectedPerm {
				t.Errorf("expected permission 0755, got %04o", info.Mode().Perm())
			}
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

func TestManager_Exists(t *testing.T) {
	tempHomeDir := filepath.Join(t.TempDir(), ".teamtime")
	tempDir := filepath.Join(tempHomeDir, ".teamtime")
	testConfigPath := filepath.Join(tempDir, "colleagues.json")

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Errorf("failed create directory %s: %v", tempDir, err)
	}
	err := os.WriteFile(testConfigPath, []byte("[]"), 0644)

	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	m, err := NewManager(tempHomeDir)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !m.Exists() {
		t.Errorf("expected file to exists")
	}
}

func TestManager_Save(t *testing.T) {
	tests := []struct {
		name           string
		colleagues     []struct{ name, city, tz string }
		wantErr        bool
		expectedLength int
	}{
		{
			name:           "empty colleague list",
			colleagues:     []struct{ name, city, tz string }{},
			wantErr:        false,
			expectedLength: 0,
		},
		{
			name: "single colleague",
			colleagues: []struct{ name, city, tz string }{
				{"Alice", "London", "Europe/London"},
			},
			wantErr:        false,
			expectedLength: 1,
		},
		{
			name: "multiple colleagues",
			colleagues: []struct{ name, city, tz string }{
				{"Alice", "London", "Europe/London"},
				{"Bob", "NYC", "America/New_York"},
				{"Charlie", "Tokyo", "Asia/Tokyo"},
			},
			wantErr:        false,
			expectedLength: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "colleagues.json")

			m := &Manager{filePath: testFile}
			cl := types.NewColleagues()

			for _, c := range tt.colleagues {
				err := cl.Add(c.name, c.city, c.tz)
				if err != nil {
					t.Fatalf("failed to add colleague: %v", err)
				}
			}

			err := m.Save(cl)

			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify file was created
			if _, err := os.Stat(testFile); err != nil {
				t.Fatalf("expected file to exist: %v", err)
			}

			// Verify content
			data, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			var loaded types.ColleagueList
			err = json.Unmarshal(data, &loaded)
			if err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			if len(loaded) != tt.expectedLength {
				t.Errorf("expected %d colleagues, got %d", tt.expectedLength, len(loaded))
			}
		})
	}

	t.Run("overwrites existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "colleagues.json")

		m := &Manager{filePath: testFile}

		// Save initial data
		cl1 := types.NewColleagues()
		cl1.Add("Alice", "London", "Europe/London")
		err := m.Save(cl1)
		if err != nil {
			t.Fatalf("first save failed: %v", err)
		}

		// Save new data (overwrites)
		cl2 := types.NewColleagues()
		cl2.Add("Bob", "NYC", "America/New_York")
		err = m.Save(cl2)
		if err != nil {
			t.Fatalf("second save failed: %v", err)
		}

		// Load and verify
		data, _ := os.ReadFile(testFile)
		var loaded types.ColleagueList
		json.Unmarshal(data, &loaded)

		if len(loaded) != 1 {
			t.Errorf("expected 1 colleague, got %d", len(loaded))
		}

		if len(loaded) > 0 && loaded[0].Name != "Bob" {
			t.Errorf("expected Bob, got %s", loaded[0].Name)
		}
	})

	t.Run("returns error for invalid path", func(t *testing.T) {
		m := &Manager{filePath: "/nonexistent/path/colleagues.json"}
		cl := types.NewColleagues()

		err := m.Save(cl)
		if err == nil {
			t.Error("expected error for invalid path")
		}
	})
}

func TestManager_Load(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		setupFile   bool
		wantCount   int
		wantErr     bool
	}{
		{
			name:        "empty array",
			fileContent: "[]",
			setupFile:   true,
			wantCount:   0,
			wantErr:     false,
		},
		{
			name:        "single colleague",
			fileContent: `[{"name":"Alice","city":"London","timezone":"Europe/London"}]`,
			setupFile:   true,
			wantCount:   1,
			wantErr:     false,
		},
		{
			name:        "multiple colleagues",
			fileContent: `[{"name":"Alice","city":"London","timezone":"Europe/London"},{"name":"Bob","city":"NYC","timezone":"America/New_York"}]`,
			setupFile:   true,
			wantCount:   2,
			wantErr:     false,
		},
		{
			name:        "empty file (zero bytes)",
			fileContent: "",
			setupFile:   true,
			wantCount:   0,
			wantErr:     false,
		},
		{
			name:        "non-existent file",
			fileContent: "",
			setupFile:   false,
			wantCount:   0,
			wantErr:     false,
		},
		{
			name:        "malformed JSON - missing bracket",
			fileContent: `[{"name":"Alice"}`,
			setupFile:   true,
			wantCount:   0,
			wantErr:     true,
		},
		{
			name:        "malformed JSON - not an array",
			fileContent: `{"name":"Alice"}`,
			setupFile:   true,
			wantCount:   0,
			wantErr:     true,
		},
		{
			name:        "invalid JSON",
			fileContent: `not json at all`,
			setupFile:   true,
			wantCount:   0,
			wantErr:     true,
		},
		{
			name:        "invalid colleague missing city and timezone",
			fileContent: `[{"name":"Alice"}]`,
			setupFile:   true,
			wantCount:   0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "colleagues.json")

			if tt.setupFile {
				err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			m := &Manager{filePath: testFile}

			cl, err := m.Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if cl != nil {
					t.Error("expected nil ColleagueList on error")
				}
				return
			}

			if cl == nil {
				t.Fatal("expected non-nil ColleagueList")
			}

			if len(*cl) != tt.wantCount {
				t.Errorf("expected %d colleagues, got %d", tt.wantCount, len(*cl))
			}

			// Additional checks for specific test cases
			if tt.name == "single colleague" && len(*cl) > 0 {
				colleagues := *cl
				if colleagues[0].Name != "Alice" {
					t.Errorf("expected colleague name Alice, got %s", colleagues[0].Name)
				}
			}

			if tt.name == "multiple colleagues" && len(*cl) == 2 {
				colleagues := *cl
				if colleagues[0].Name != "Alice" {
					t.Errorf("expected first colleague Alice, got %s", colleagues[0].Name)
				}
				if colleagues[1].Name != "Bob" {
					t.Errorf("expected second colleague Bob, got %s", colleagues[1].Name)
				}
			}

			if tt.name == "file with missing fields" && len(*cl) > 0 {
				colleagues := *cl
				if colleagues[0].Name != "Alice" {
					t.Errorf("expected name Alice, got %s", colleagues[0].Name)
				}
				if colleagues[0].City != "" {
					t.Errorf("expected empty city, got %s", colleagues[0].City)
				}
			}
		})
	}
}

func TestManager_Integration(t *testing.T) {
	tests := []struct {
		name       string
		colleagues []struct{ name, city, tz string }
	}{
		{
			name:       "empty list",
			colleagues: []struct{ name, city, tz string }{},
		},
		{
			name: "single colleague",
			colleagues: []struct{ name, city, tz string }{
				{"Alice", "London", "Europe/London"},
			},
		},
		{
			name: "multiple colleagues",
			colleagues: []struct{ name, city, tz string }{
				{"Alice", "London", "Europe/London"},
				{"Bob", "NYC", "America/New_York"},
				{"Charlie", "Tokyo", "Asia/Tokyo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "colleagues.json")

			m := &Manager{filePath: testFile}

			// Create and save data
			original := types.NewColleagues()
			for _, c := range tt.colleagues {
				err := original.Add(c.name, c.city, c.tz)
				if err != nil {
					t.Fatalf("failed to add colleague: %v", err)
				}
			}

			err := m.Save(original)
			if err != nil {
				t.Fatalf("save failed: %v", err)
			}

			// Load into new list
			loaded, err := m.Load()
			if err != nil {
				t.Fatalf("load failed: %v", err)
			}

			if loaded == nil {
				t.Fatal("expected non-nil ColleagueList")
			}

			// Compare
			if len(*loaded) != len(*original) {
				t.Fatalf("expected %d colleagues, got %d", len(*original), len(*loaded))
			}

			origList := *original
			loadedList := *loaded

			for i := range origList {
				if origList[i].Name != loadedList[i].Name {
					t.Errorf("colleague %d: expected name %s, got %s", i, origList[i].Name, loadedList[i].Name)
				}
				if origList[i].City != loadedList[i].City {
					t.Errorf("colleague %d: expected city %s, got %s", i, origList[i].City, loadedList[i].City)
				}
				if origList[i].Timezone != loadedList[i].Timezone {
					t.Errorf("colleague %d: expected timezone %s, got %s", i, origList[i].Timezone, loadedList[i].Timezone)
				}
			}
		})
	}
}
