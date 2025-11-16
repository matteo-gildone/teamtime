package cmd

import (
	"fmt"
	"strconv"

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
	m, _ := GetManager(cmd.Context())
	colleagues, _ := GetColleagues(cmd.Context())
	err = colleagues.Delete(idx)

	if err != nil {
		return fmt.Errorf("failed remove 'colleagues.json' in: %w", err)
	}

	err = m.Save(colleagues)

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
