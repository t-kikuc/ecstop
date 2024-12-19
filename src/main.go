package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/src/stop"
)

func main() {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "ecstop",
		Short: "Stop ECS resources easily",
		Long:  ``,
	}

	rootCmd.AddCommand(
		stop.NewStopServiceCommand(),
		stop.NewStopTaskCommand(),
		stop.NewStopInstanceCommand(),
		stop.NewStopAllCommand(),
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
