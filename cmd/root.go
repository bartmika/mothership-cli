package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	port        int
	databaseUrl string
	hmacSecret  string
)

var rootCmd = &cobra.Command{
	Use:   "mothership-cli",
	Short: "Command line interface for the mothership server",
	Long:  `The purpose of this application is to provide sub-commands for users to perform actions on their account in remote mothership server.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
