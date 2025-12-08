package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/matteo-gildone/teamtime/internals/config"
	"github.com/matteo-gildone/teamtime/internals/types"
	"github.com/spf13/cobra"
)

type contextKey string

const (
	managerKey    contextKey = "manager"
	colleaguesKey contextKey = "colleagues"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "teamtime",
	Short:        "Global team time & weather",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" || cmd.Name() == "help" {
			return nil
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory %w", err)
		}
		m, err := config.NewManager(homeDir)

		if err != nil {
			return fmt.Errorf("failed to create manager %w", err)
		}

		if !m.Exists() {
			return fmt.Errorf("'colleagues.json' not found, run 'teamtime init'")
		}

		cl, err := m.Load()
		if err != nil {
			return fmt.Errorf("failed load 'colleagues.json': %w", err)
		}

		ctx := cmd.Context()
		ctx = context.WithValue(ctx, managerKey, m)
		ctx = context.WithValue(ctx, colleaguesKey, cl)

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

func GetColleagues(ctx context.Context) (*types.ColleagueList, error) {
	val := ctx.Value(colleaguesKey)
	colleagues, ok := val.(*types.ColleagueList)
	if !ok {
		return nil, fmt.Errorf("%s not found in context", colleaguesKey)
	}
	return colleagues, nil
}

func GetManager(ctx context.Context) (*config.Manager, error) {
	val := ctx.Value(managerKey)
	m, ok := val.(*config.Manager)
	if !ok || m == nil {
		return nil, fmt.Errorf("%s not found in context", managerKey)
	}
	return m, nil
}
