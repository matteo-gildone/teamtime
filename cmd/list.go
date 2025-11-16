package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/matteo-gildone/teamtime/internals/config"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show current local time for all colleagues",
	RunE:  listFunc,
}

func listFunc(cmd *cobra.Command, args []string) error {
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

	now := time.Now()

	fmt.Printf("| %-20s | %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	fmt.Printf("| %-20s | %-20s | %-20s | %-20s | %-20s |\n", "ID", "Name", "City", "Timezone", "Local Time")
	fmt.Printf("| %-20s | %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	for idx, c := range *cl {
		loc, _ := time.LoadLocation(c.Timezone)
		local := now.In(loc)
		fmt.Printf("| %-20d | %-20s | %-20s | %-20s | %-20s |\n", idx+1, c.Name, c.City, c.Timezone, local.Format("15:04 (Mon 02 Jan)"))
	}

	fmt.Printf("| %-20s | %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
