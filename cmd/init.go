package cmd

import (
	"fmt"
	"os"

	"github.com/matteo-gildone/teamtime/internals/storage"
	"github.com/matteo-gildone/teamtime/internals/styles"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
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
	m, err := storage.NewManager(homeDir)

	if err != nil {
		return fmt.Errorf("init command - %w", err)
	}

	if m.Exists() {
		return fmt.Errorf("app already initialised app in %s\n", m.GetRelativeFilePath())
	}

	if err = m.EnsureFolder(); err != nil {
		return fmt.Errorf("ensure folder - %w", err)
	}

	cl := types.ColleagueList{}

	if err = m.Save(&cl); err != nil {
		return fmt.Errorf("failed create 'colleagues.json' %w", err)
	}

	fmt.Println(styles.NewStyles().Cyan().Bold().Render(`
 ____  ____   __   _  _  ____  __  _  _  ____
(_  _)(  __) / _\ ( \/ )(_  _)(  )( \/ )(  __)
  )(   ) _) /    \/ \/ \  )(   )( / \/ \ ) _)
 (__) (____)\_/\_/\_)(_/ (__) (__)\_)(_/(____)`))

	successStyle := styles.NewStyles().Green()
	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf("Initialised app in: %s", m.GetRelativeFilePath())))
	fmt.Println()
	return nil
}
