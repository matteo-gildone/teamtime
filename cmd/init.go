package cmd

import (
	"fmt"
	"os"

	"github.com/matteo-gildone/teamtime/internals/config"
	"github.com/matteo-gildone/teamtime/internals/styles"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		err := initFunc()
		if err != nil {
			return fmt.Errorf("init command:%w", err)
		}
		return nil
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
	m, err := config.NewManager(homeDir)

	if err != nil {
		return fmt.Errorf("init command - %w", err)
	}

	if m.Exists() {
		return fmt.Errorf("app already initialised app in %s\n", m.GetFilePath())
	}

	if err = m.EnsureFolder(); err != nil {
		return fmt.Errorf("ensure folder - %w", err)
	}

	cl := types.ColleagueList{}

	if err = m.Save(&cl); err != nil {
		return fmt.Errorf("failed create 'colleagues.json' %w", err)
	}

	successStyle := styles.NewStyles().Green()
	successMessage := fmt.Sprintf("Initialised app in: %s\n", m.GetFilePath())
	styledSuccessMessage := successStyle.Render(successMessage)
	fmt.Println()
	fmt.Println(styledSuccessMessage)
	return nil
}
