package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

// checkCmd represents the list command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Show current local time for all colleagues",
	Args:  cobra.ExactArgs(1),
	RunE:  checkFunc,
}

func checkFunc(cmd *cobra.Command, args []string) error {
	colleagues, err := GetColleagues(cmd.Context())

	if err != nil {
		return err
	}

	if len(*colleagues) == 0 {
		fmt.Println("no colleague present")
		return nil
	}

	if args[0] == "all" {
		renderTable(*colleagues)
		return nil
	}

	var filteredColleagues types.ColleagueList
	for _, c := range *colleagues {
		if strings.EqualFold(c.Name, args[0]) {
			filteredColleagues = append(filteredColleagues, c)
		}
	}

	if len(filteredColleagues) == 0 {
		fmt.Printf("no colleague found with name: %s\n", args[0])
		return nil
	}

	renderTable(filteredColleagues)
	return nil

}

func renderTable(colleagues types.ColleagueList) {
	if len(colleagues) == 0 {
		return
	}
	now := time.Now()

	fmt.Println()
	fmt.Printf("%-20s | %-20s | %-20s | %-20s\n", "ID", "Name", "City", "Local Time")
	fmt.Printf("%-20s | %-20s | %-20s | %-20s\n", strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20), strings.Repeat("-", 20))
	for idx, c := range colleagues {
		loc, err := time.LoadLocation(c.Timezone)
		if err != nil {
			fmt.Printf("%-20d | %-20s | %-20s | %-20s\n", idx+1, c.Name, c.City, "ERROR: Invalid TZ")
			continue
		}
		local := now.In(loc)
		fmt.Printf("%-20d | %-20s | %-20s | %-20s\n", idx+1, c.Name, c.City, local.Format("15:04 (Mon 02 Jan)"))
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
