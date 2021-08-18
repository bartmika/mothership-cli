package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// The following are required.
	relayCmd.Flags().StringVarP(&iAccessToken, "access_token", "a", "", "The JWT access token provided with successful authentication")
	relayCmd.MarkFlagRequired("access_token")
	relayCmd.Flags().StringVarP(&iRefreshToken, "refresh_token", "b", "", "The JWT refresh token provided with successful authentication")
	relayCmd.MarkFlagRequired("refresh_token")

	// The following are optional and will have defaults placed when missing.
	relayCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(relayCmd)
}

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Relay local data to remote",
	Long:  `Connect to remote mothership server and serve a local gRPC server. Local apps can submit data this relay through the local gRPC server. Relay will take all local gRPC server submissions and send them to the remote mothership server for the authenticated credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: IMPL.")
	},
}
