package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecstop/src/stop"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecstop",
	Short: "Stop ECS resources",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	// Register commands
	commands := []*cobra.Command{
		stop.NewStopServiceCommand(),
		stop.NewStopTaskCommand(),
		stop.NewStopInstanceCommand(),
		stop.NewStopAllCommand(),
		// deletetaskdef.NewCommand(),
		// deletetasksets.NewCommand(),
	}
	for _, c := range commands {
		rootCmd.AddCommand(c)
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
