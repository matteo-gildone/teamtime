package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matteo-gildone/teamtime/internals/types"
)

const maxFileSize = 10 * 1024 * 1024 // 10MB

var ErrMissingHomeDir = errors.New("'homeDir' must not be empty")

type Manager struct {
	homeDir  string
	filePath string
}

func (m *Manager) Save(cl *types.ColleagueList) error {
	js, err := json.Marshal(cl)
	if err != nil {
		return err
	}

	return os.WriteFile(m.filePath, js, 0600)
}

func (m *Manager) Load() (*types.ColleagueList, error) {
	info, err := os.Stat(m.filePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return types.NewColleagues(), nil
	}

	if info != nil && info.Size() > maxFileSize {
		return nil, fmt.Errorf("file too large %d bytes (max %d)", info.Size(), maxFileSize)
	}

	file, err := os.ReadFile(m.filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return types.NewColleagues(), nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(file) == 0 {
		return types.NewColleagues(), nil
	}

	cl := types.NewColleagues()
	if err := json.Unmarshal(file, cl); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for i, c := range *cl {
		if err = c.Validate(); err != nil {
			return nil, fmt.Errorf("colleague at index %d: %w", i+1, err)
		}
	}

	return cl, nil
}
func (m *Manager) Exists() bool {
	_, err := os.Stat(m.filePath)
	return err == nil
}

func (m *Manager) EnsureFolder() error {
	configDir := filepath.Join(m.homeDir, ".teamtime")

	if _, err := os.Stat(configDir); err == nil {
		return nil
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed create directory %s: %w", configDir, err)
	}
	return nil
}

func (m *Manager) GetFilePath() string {
	return m.filePath
}

func (m *Manager) GetRelativeFilePath() string {
	rel, err := filepath.Rel(m.homeDir, m.filePath)
	if err != nil {
		return m.filePath
	}
	return filepath.Join("~", rel)
}

func NewManager(homeDir string) (*Manager, error) {
	if homeDir == "" {
		return nil, ErrMissingHomeDir
	}
	cleanHome := filepath.Clean(homeDir)
	if !filepath.IsAbs(cleanHome) {
		return nil, fmt.Errorf("homeDir must be an absolute path")
	}
	configPath := filepath.Join(homeDir, ".teamtime", "colleagues.json")

	return &Manager{
		homeDir:  homeDir,
		filePath: configPath,
	}, nil
}
