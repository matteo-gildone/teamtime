package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/matteo-gildone/teamtime/internals/config"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [name] [city] [time zone]",
	Short: "Add a new colleague",
	Args:  cobra.ExactArgs(3),
	RunE:  addFunc,
}

func addFunc(cmd *cobra.Command, args []string) error {
	if _, err := time.LoadLocation(args[2]); err != nil {
		return fmt.Errorf("invilide timezone: %s %w", args[2], err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory %w", err)
	}
	configPath := filepath.Join(homeDir, ".teamtime", "colleagues.json")
	m := config.NewManager(configPath)
	cl := types.NewColleagues()
	if err := m.Load(cl); err != nil {
		return fmt.Errorf("failed load 'colleagues.json' in: %w", err)
	}
	cl.Add(args[0], args[1], args[2])
	fmt.Println(configPath)
	err = m.Save(cl)

	if err != nil {
		return fmt.Errorf("failed add 'colleagues.json' in: %w", err)
	}
	fmt.Printf("%s was added\n", args[0])
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
