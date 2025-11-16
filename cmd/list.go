package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show current local time for all colleagues",
	RunE:  listFunc,
}

func listFunc(cmd *cobra.Command, args []string) error {
	colleagues, err := GetColleagues(cmd.Context())

	if err != nil {
		return err
	}

	if len(*colleagues) == 0 {
		fmt.Println("no colleague present")
		return nil
	}

	now := time.Now()

	fmt.Printf("| %-20s | %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	fmt.Printf("| %-20s | %-20s | %-20s | %-20s | %-20s |\n", "ID", "Name", "City", "Timezone", "Local Time")
	fmt.Printf("| %-20s | %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	for idx, c := range *colleagues {
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
