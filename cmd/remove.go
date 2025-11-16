package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/matteo-gildone/teamtime/internals/config"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove colleague",
	Args:  cobra.ExactArgs(1),
	RunE:  removeFunc,
}

func removeFunc(cmd *cobra.Command, args []string) error {
	idx, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("idx must be a number %w", err)
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

	err = cl.Delete(idx)

	if err != nil {
		return fmt.Errorf("failed remove 'colleagues.json' in: %w", err)
	}

	err = m.Save(cl)

	if err != nil {
		return fmt.Errorf("failed add 'colleagues.json' in: %w", err)
	}
	fmt.Printf("%s was removed\n", args[0])
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
