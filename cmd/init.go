package cmd

import (
	"fmt"

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
	m, err := config.NewManager()

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

	fmt.Printf("Initialise app in: %s\n", m.GetFilePath())
	return nil
}
