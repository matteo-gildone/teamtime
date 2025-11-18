package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/matteo-gildone/teamtime/internals/types"
)

type Manager struct {
	filePath string
}

func (m *Manager) Save(cl *types.ColleagueList) error {
	js, err := json.Marshal(cl)
	if err != nil {
		return err
	}

	return os.WriteFile(m.filePath, js, 0664)
}

func (m *Manager) Load(cl *types.ColleagueList) error {
	file, err := os.ReadFile(m.filePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, cl)
}
func (m *Manager) Exists() bool {
	_, err := os.Stat(m.filePath)
	return err == nil
}

func (m *Manager) EnsureFolder() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory %w", err)
	}
	configDir := filepath.Join(homeDir, ".teamtime")

	if _, err = os.Stat(configDir); err == nil {
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

func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory %w", err)
	}
	configPath := filepath.Join(homeDir, ".teamtime", "colleagues.json")

	return &Manager{
		filePath: configPath,
	}, nil
}
