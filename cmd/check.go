package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/matteo-gildone/teamtime/internals/service"
	"github.com/matteo-gildone/teamtime/internals/styles"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

const (
	extendedStart  int = 7
	extendedEnd    int = 20
	workHoursStart int = 9
	workHoursEnd   int = 17
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

func init() {
	checkCmd.Flags().BoolP("watch", "w", false, "continuously update times")
	checkCmd.Flags().IntP("interval", "i", 10, "update interval ins minutes")
	rootCmd.AddCommand(checkCmd)
}

func checkFunc(cmd *cobra.Command, args []string) error {
	svc, err := GetColleaguesService(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get colleague service: %w", err)
	}

	watchMode, err := cmd.Flags().GetBool("watch")
	if err != nil {
		return fmt.Errorf("failed to get watch flag: %w", err)
	}

	if watchMode {
		watchInterval, err := cmd.Flags().GetInt("interval")
		if err != nil {
			return fmt.Errorf("failed to get interval flag: %w", err)
		}
		return runWatch(cmd.Context(), svc, args[0], watchInterval)
	}

	return runOnce(svc, args[0])
}

func runOnce(svc *service.ColleagueService, query string) error {
	colleagues, err := getColleagues(svc, query)
	if err != nil {
		return err
	}

	displayColleagues(colleagues, query)
	return nil
}

func runWatch(ctx context.Context, svc *service.ColleagueService, query string, interval int) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	if err := renderWatchScreen(svc, query, interval); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			msgStyle := styles.NewStyles().Cyan()
			clearScreen()
			fmt.Println(msgStyle.Render("exiting watch mode..."))
			return nil
		case <-ticker.C:
			if err := renderWatchScreen(svc, query, interval); err != nil {
				return err
			}

		}
	}

}

func getColleagues(svc *service.ColleagueService, query string) (types.ColleagueList, error) {
	if query == "all" {
		return svc.AllColleagues()
	}

	return svc.FindColleague(query)
}

func displayColleagues(colleagues []types.Colleague, query string) {
	if len(colleagues) == 0 {
		displayEmptyMessage(query)
		return
	}

	renderTable(colleagues)
}

func displayEmptyMessage(query string) {
	msgStyle := styles.NewStyles().Cyan()
	if query == "all" {
		fmt.Println(msgStyle.Render("no colleagues found"))
	} else {
		fmt.Println(msgStyle.Render(fmt.Sprintf("no colleague found with name: %q\n", query)))
	}

}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func renderWatchScreen(svc *service.ColleagueService, query string, interval int) error {
	clearScreen()
	colleagues, err := getColleagues(svc, query)
	if err != nil {
		return err
	}

	watchStyle := styles.NewStyles().Cyan()
	dimStyle := styles.NewStyles().Dim()

	fmt.Println(watchStyle.Render(fmt.Sprintf("⟳ Watch mode (updates every %d mins) - Press Ctrl+C to exit", interval)))
	fmt.Println()
	displayColleagues(colleagues, query)
	fmt.Println(dimStyle.Render(fmt.Sprintf("Last updated: %s", time.Now().Format("15:04:05"))))
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
		timeDisplay := getDisplayTime(local, plainStyle)
		fmt.Printf("%-4d | %-20s | %s\n",
			idx+1,
			c.Name,
			timeDisplay)
	}
	fmt.Println()
	renderLegend(plainStyle)
}

func classifyTimeOfDay(hour int) timeClassification {
	if hour >= workHoursStart && hour < workHoursEnd {
		return timeWork
	}

	if (hour >= extendedStart && hour < workHoursStart) || (hour >= workHoursEnd && hour < extendedEnd) {
		return timeExtended
	}

	return timeOff
}

func getDisplayTime(localTime time.Time, plainStyle styles.Style) string {
	hour := localTime.Hour()
	timeStr := localTime.Format("15:04 (Mon 02 Jan)")
	base := plainStyle.Bold()

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

func renderLegend(plainStyle styles.Style) {
	if plainStyle.NoColor() {
		return
	}
	fmt.Println(plainStyle.Render("Availability:"))
	fmt.Println(plainStyle.Cyan().Bold().Render("    Cyan") + " - Work hours (9am-5pm)")
	fmt.Println(plainStyle.Yellow().Bold().Render("    Yellow") + " - Extended hours")
	fmt.Println(plainStyle.Red().Bold().Render("    Red") + " - Off hours")
	fmt.Println()
}
