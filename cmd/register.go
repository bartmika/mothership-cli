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
	reqEmail     string
	reqPassword  string
	reqCompany   string
	reqFirstName string
	reqLastName  string
	reqTimezone  string
)

func init() {
	// The following are required.
	registerCmd.Flags().StringVarP(&reqEmail, "email", "e", "", "The email you want to associate with your account")
	registerCmd.MarkFlagRequired("email")
	registerCmd.Flags().StringVarP(&reqPassword, "password", "x", "", "The password you want to use to protect your account")
	registerCmd.MarkFlagRequired("password")
	registerCmd.Flags().StringVarP(&reqFirstName, "first_name", "f", "", "Your first name to use in your account")
	registerCmd.MarkFlagRequired("first_name")
	registerCmd.Flags().StringVarP(&reqLastName, "last_name", "l", "", "Your last name to use in your account")
	registerCmd.MarkFlagRequired("last_name")

	// The following are optional and will have defaults placed when missing.
	registerCmd.Flags().StringVarP(&reqCompany, "company", "c", "", "Your companies name. If personal use then leave blank")
	registerCmd.Flags().StringVarP(&reqTimezone, "timezone", "t", "America/Toronto", "Your accounts timezone.`")
	registerCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of to connect to the server")

	// Setup our sub-command.
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Creates your account",
	Long:  `Connects to the mothership server and creates your account`,
	Run: func(cmd *cobra.Command, args []string) {
		doRegisterCmd()
	},
}

func doRegisterCmd() {
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

	req := &pb.RegistrationReq{
		Email:     reqEmail,
		Password:  reqPassword,
		Company:   reqCompany,
		FirstName: reqFirstName,
		LastName:  reqLastName,
		Timezone:  reqTimezone,
	}

	// Perform our gRPC request.
	res, err := client.Register(ctx, req)
	if err != nil {
		log.Fatalf("Registration failed:\n%v\n", err)
	}


	log.Println(res.Message)
}
