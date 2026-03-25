package cmd

import (

	"github.com/spf13/cobra"
)

var minutes int

var startCmd = &cobra.Command{
	Use: "start",
	Short: "Start the process",
	Run: func(cmd *cobra.Command, args []string) {
			},
}

func init() {
	startCmd.Flags().IntVarP(&minutes, "min", "m", 1, "Minutes to run")
	rootCmd.AddCommand(startCmd)
}
