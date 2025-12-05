package cmd

import (
	"fmt"

	"github.com/matteo-gildone/teamtime/internals/styles"
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
	m, err := GetManager(cmd.Context())
	if err != nil {
		return err
	}
	colleagues, err := GetColleagues(cmd.Context())
	if err != nil {
		return err
	}

	newColleague, err := types.NewColleague(args[0], args[1], args[2])

	if err != nil {
		return fmt.Errorf("invalid colleague data: %w", err)
	}

	colleagues.Add(newColleague)

	err = m.Save(colleagues)

	if err != nil {
		return fmt.Errorf("failed to save to '%s' in: %w", m.GetFilePath(), err)
	}

	successStyle := styles.NewStyles().Green()
	fmt.Println(successStyle.Render(fmt.Sprintf("✓ %s was added", args[0])))
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
