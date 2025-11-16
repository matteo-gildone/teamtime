package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matteo-gildone/teamtime/internals/config"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

const (
	appDir = ".teamtime"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise project",
	Run: func(cmd *cobra.Command, args []string) {
		err := initFunc()
		if err != nil {
			fmt.Println(fmt.Errorf("init command:%w", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initFunc() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory %w", err)
	}
	configDir := filepath.Join(homeDir, appDir)
	_, err = os.Stat(configDir)
	if err == nil {
		return fmt.Errorf("Reinitialized existing app in %s\n\n", configDir)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed create directory %s: %w", configDir, err)
	}
	configPath := filepath.Join(configDir, "colleagues.json")

	m := config.NewManager(configPath)

	cl := types.ColleagueList{}

	err = m.Save(&cl)

	if err != nil {
		return fmt.Errorf("failed create 'colleagues.json' in %s: %w", configDir, err)
	}

	fmt.Printf("Initialise app in: %s\n", configDir)
	return nil
}
