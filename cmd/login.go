package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	// tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/bartmika/mothership-server/proto"
)

var (
	logEmail    string
	logPassword string
)

func init() {
	// The following are required.
	loginCmd.Flags().StringVarP(&logEmail, "email", "e", "", "The email you want to associate with your account")
	loginCmd.MarkFlagRequired("email")
	loginCmd.Flags().StringVarP(&logPassword, "password", "x", "", "The password you want to use to protect your account")
	loginCmd.MarkFlagRequired("password")

	// The following are optional and will have defaults placed when missing.
	loginCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of to connect to the server")

	// Setup our sub-command.
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into your account",
	Long:  `Connects to the mothership server and log into your account`,
	Run: func(cmd *cobra.Command, args []string) {
		doLoginCmd()
	},
}

func doLoginCmd() {
	// Set up a direct connection to the gRPC server.
	conn, err := grpc.Dial(
		fmt.Sprintf(":%v", port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Set up our protocol buffer interface.
	client := pb.NewMothershipClient(conn)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.LoginReq{
		Email:    logEmail,
		Password: logPassword,
	}

	// Perform our gRPC request.
	res, err := client.Login(ctx, req)
	if err != nil {
		log.Fatalf("Registration failed:\n%v\n", err)
	}

	log.Printf("Successfully logged in. Your credentials are as follows, please run this in your console.\n\n")
	fmt.Println("Access Token:")
	fmt.Printf("export MOTHERSHIP_CLI_ACCESS_TOKEN=%v\n\n",res.AccessToken)
	fmt.Println("Refresh Token:")
	fmt.Printf("export MOTHERSHIP_CLI_REFRESH_TOKEN=%v\n\n", res.RefreshToken)
}
