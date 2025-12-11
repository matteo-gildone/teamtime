package cmd

import (
	"fmt"
	"strconv"

	"github.com/matteo-gildone/teamtime/internals/styles"
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
		return fmt.Errorf("index must be a number %w", err)
	}
	svc, err := GetColleaguesService(cmd.Context())
	if err != nil {
		return err
	}

	removed, err := svc.RemoveColleague(idx)
	if err != nil {
		return fmt.Errorf("remove command: %w", err)
	}

	successStyle := styles.NewStyles().Green()
	fmt.Println(successStyle.Render(fmt.Sprintf("✓ %s was removed", removed.Name)))
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
