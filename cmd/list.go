package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show current local time for all colleagues",
	Run:   listFunc,
}

func listFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("| %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	fmt.Printf("| %-20s | %-20s | %-20s | %-20s |\n", "Name", "City", "Timezone", "Local Time")
	fmt.Printf("| %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	fmt.Printf("| %-20s | %-20s | %-20s | %-20s |\n", "Matteo", "Verona", "Europe/Italy", "10:30")
	fmt.Printf("| %-20s | %-20s | %-20s | %-20s |\n", "Cookin Monster", "Tuxon", "US/Arizona", "20:30")
	fmt.Printf("| %-20s + %-20s + %-20s + %-20s |\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
