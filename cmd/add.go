package cmd

import (
	"fmt"

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
	err = colleagues.Add(args[0], args[1], args[2])

	if err != nil {
		return fmt.Errorf("failed add colleague: %w", err)
	}

	err = m.Save(colleagues)

	if err != nil {
		return fmt.Errorf("failed save 'colleagues.json' in: %w", err)
	}

	fmt.Printf("%s was added\n", args[0])
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
