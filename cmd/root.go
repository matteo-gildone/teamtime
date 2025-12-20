package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/matteo-gildone/teamtime/internals/service"
	"github.com/matteo-gildone/teamtime/internals/storage"
	"github.com/spf13/cobra"
)

type contextKey string

const (
	serviceKey contextKey = "colleagueservice"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "teamtime",
	Short:        "Global team time",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" || cmd.Name() == "help" {
			return nil
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory %w", err)
		}
		m, err := storage.NewManager(homeDir)

		if err != nil {
			return fmt.Errorf("failed to create manager %w", err)
		}

		if !m.Exists() {
			return fmt.Errorf("'colleagues.json' not found, run 'teamtime init'")
		}

		svc := service.NewColleagueService(m)

		ctx := cmd.Context()
		ctx = context.WithValue(ctx, serviceKey, svc)

		cmd.SetContext(ctx)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func GetColleaguesService(ctx context.Context) (*service.ColleagueService, error) {
	val := ctx.Value(serviceKey)
	if val == nil {
		return nil, fmt.Errorf("service not found in context")
	}
	svc, ok := val.(*service.ColleagueService)
	if !ok {
		return nil, fmt.Errorf("unexpected type in context, got: %T, want *service.ColleagueService", val)
	}
	return svc, nil
}
