package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/t-kikuc/ecscale0/src/scalein"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecscale0",
	Short: "Scale-in ECS Services and Tasks",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	// Register commands
	commands := []*cobra.Command{
		scalein.NewScaleinServiceCommand(),
		scalein.NewStopTaskCommand(),
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
