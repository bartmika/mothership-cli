package cmd

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/metadata"

	pb "github.com/bartmika/mothership-server/proto"
)

// Ex:
// go run main.go insert -p=50051 -m="solar_biodigester_temperature_in_degrees" -v=50 -t=1600000000 -a=xxx -b=yyy

var (
	metric        string
	value         float64
	tsv           int64
	iAccessToken  string
	iRefreshToken string
)

func init() {
	// The following are required.
	insertRowCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("metric")
	insertRowCmd.Flags().Float64VarP(&value, "value", "v", 0.00, "The value to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("value")
	insertRowCmd.Flags().Int64VarP(&tsv, "timestamp", "t", 0, "The timestamp to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("timestamp")

	// The following are optional and will have defaults placed when missing.
	insertRowCmd.Flags().StringVarP(&iAccessToken, "access_token", "a", os.Getenv("MOTHERSHIP_CLI_ACCESS_TOKEN"), "The JWT access token for your account. Leave blank to access environment variable.")
	insertRowCmd.Flags().StringVarP(&iRefreshToken, "refresh_token", "b", os.Getenv("MOTHERSHIP_CLI_REFRESH_TOKEN"), "The JWT refresh token for your account. Leave blank to access environment variable.")
	insertRowCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(insertRowCmd)
}

func doInsertRow() {
	// Here is the code which attaches our authorization information to our
	// context and has the context sent to the server with these credentials.
	// Use this context ONLY when making RPC calls.
	//
	// Special Thanks:
	// https://shijuvar.medium.com/writing-grpc-interceptors-in-go-bf3e7671fe48
	ctx := context.Background()
	md := metadata.Pairs("authorization", iAccessToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

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

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	ts := &tspb.Timestamp{
		Seconds: tsv,
		Nanos:   0,
	}

	// Generate our labels.
	labels := []*pb.Label{}
	labels = append(labels, &pb.Label{Name: "Source", Value: "Command"})

	// Perform our gRPC request.
	_, err = client.InsertTimeSeriesDatum(ctx, &pb.TimeSeriesDatumReq{Labels: labels, Metric: metric, Value: value, Timestamp: ts})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}

	log.Printf("Successfully inserted")
}

var insertRowCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a single time-series datum",
	Long:  `Connect to the gRPC server and send a single time-series datum.`,
	Run: func(cmd *cobra.Command, args []string) {
		if iAccessToken == "" {
			at := os.Getenv("MOTHERSHIP_CLI_ACCESS_TOKEN")
			if at != "" {
				iAccessToken = at
			} else {
				log.Fatal("No access token set.")
			}
		}
		if iRefreshToken == "" {
			rt := os.Getenv("MOTHERSHIP_CLI_REFRESH_TOKEN")
			if rt != "" {
				iRefreshToken = rt
			} else {
				log.Fatal("No refresh token set.")
			}
		}
		doInsertRow()
	},
}
