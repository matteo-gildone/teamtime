package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/matteo-gildone/teamtime/internals/styles"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

type timeClassification string

const (
	timeWork     timeClassification = "work"
	timeExtended timeClassification = "extended"
	timeOff      timeClassification = "off"
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
		fmt.Println("no colleague found")
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
		fmt.Printf("no colleague found with name: %q\n", args[0])
		return nil
	}

	renderTable(filteredColleagues)
	return nil

}

func renderTable(colleagues types.ColleagueList) {
	plainStyle := styles.NewStyles()
	heading := plainStyle.Bold()
	invalidTZ := heading.Red()
	if len(colleagues) == 0 {
		return
	}
	now := time.Now()

	fmt.Println()
	fmt.Printf("%s | %s | %s\n",
		heading.Render(fmt.Sprintf("%-4s", "ID")),
		heading.Render(fmt.Sprintf("%-20s", "Name")),
		heading.Render(fmt.Sprintf("%-32s", "Local Time")))

	fmt.Printf("%-4s | %-20s | %-20s\n",
		strings.Repeat("-", 4),
		strings.Repeat("-", 20),
		strings.Repeat("-", 32))
	for idx, c := range colleagues {
		loc, err := time.LoadLocation(c.Timezone)
		if err != nil {
			fmt.Printf("%-4d | %-20s | %s\n",
				idx+1,
				c.Name,
				invalidTZ.Render(fmt.Sprintf("%-32s", "ERROR: Invalid TZ")))
			continue
		}
		local := now.In(loc)
		timeDisplay := getDisplayTime(local)
		fmt.Printf("%-4d | %-20s | %s\n",
			idx+1,
			c.Name,
			timeDisplay)
	}
	fmt.Println()
	if !styles.NoColor {
		renderLegend()
	}
}

func classifyTimeOfDay(hour int) timeClassification {
	if hour >= 9 && hour < 17 {
		return timeWork
	}

	if (hour >= 7 && hour < 9) || (hour >= 17 && hour < 19) {
		return timeExtended
	}

	return timeOff
}

func getDisplayTime(localTime time.Time) string {
	hour := localTime.Hour()
	timeStr := localTime.Format("15:04 (Mon 02 Jan)")
	base := styles.NewStyles().Bold()

	switch classifyTimeOfDay(hour) {
	case timeWork:
		return base.Cyan().Render(fmt.Sprintf("%-32s", timeStr))
	case timeExtended:
		return base.Yellow().Render(fmt.Sprintf("%-32s", timeStr+" [Extended]"))
	case timeOff:
		return base.Red().Render(fmt.Sprintf("%-32s", timeStr+" [Off]"))
	default:
		return base.Render(fmt.Sprintf("%-32s", timeStr))
	}
}

func renderLegend() {
	plain := styles.NewStyles()
	fmt.Println(plain.Render("Availability:"))
	fmt.Println(plain.Cyan().Bold().Render("    Cyan") + " - Work hours (9am-5pm)")
	fmt.Println(plain.Yellow().Bold().Render("    Yellow") + " - Extended hours")
	fmt.Println(plain.Red().Bold().Render("    Red") + " - Off hours")
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
