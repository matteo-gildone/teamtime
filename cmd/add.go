package cmd

import (
	"fmt"

	"github.com/matteo-gildone/teamtime/internals/styles"
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
	svc, err := GetColleaguesService(cmd.Context())
	if err != nil {
		return err
	}

	newColleague, err := svc.AddColleague(args[0], args[1], args[2])
	if err != nil {
		return fmt.Errorf("add command: %w", err)
	}

	successStyle := styles.NewStyles().Green()
	fmt.Println(successStyle.Render(fmt.Sprintf("✓ %s was added", newColleague.Name)))
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
