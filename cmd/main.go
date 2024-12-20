package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/pkg/stop"
)

func main() {
	if err := executeCmd(); err != nil {
		os.Exit(1)
	}
}

func executeCmd() error {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "ecstop",
		Short: "ecstop stops ECS resources instantly",
		Long:  ``,
	}

	rootCmd.AddCommand(
		stop.NewStopServiceCommand(),
		stop.NewStopTaskCommand(),
		stop.NewStopInstanceCommand(),
		stop.NewStopAllCommand(),
	)

	return rootCmd.Execute()
}
