package config

import (
	"encoding/json"
	"errors"
	"os"

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

func NewManager(path string) *Manager {
	return &Manager{
		filePath: path,
	}
}
