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
// go run main.go select --port=50051 --metric="bio_reactor_pressure_in_kpa" --start=1600000000 --end=1725946120 -a=xxx -b=yyy

var (
	start int64
	end   int64
)

func init() {
	// The following are required.
	selectBulkCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to filter by")
	selectBulkCmd.MarkFlagRequired("metric")
	selectBulkCmd.Flags().Int64VarP(&start, "start", "s", 0, "The start timestamp to begin our range")
	selectBulkCmd.MarkFlagRequired("start")
	selectBulkCmd.Flags().Int64VarP(&end, "end", "e", 0, "The end timestamp to finish our range")
	selectBulkCmd.MarkFlagRequired("end")

	// The following are optional and will have defaults placed when missing.
	selectBulkCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	selectBulkCmd.Flags().StringVarP(&iAccessToken, "access_token", "a", os.Getenv("MOTHERSHIP_CLI_ACCESS_TOKEN"), "The JWT access token for your account. Leave blank to access environment variable.")
	selectBulkCmd.Flags().StringVarP(&iRefreshToken, "refresh_token", "b", os.Getenv("MOTHERSHIP_CLI_REFRESH_TOKEN"), "The JWT refresh token for your account. Leave blank to access environment variable.")
	rootCmd.AddCommand(selectBulkCmd)
}

func doSelectRow() {
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

	ctx, cancel := context.WithTimeout(ctx, 60 * time.Second)
	defer cancel()

	// Convert the unix timestamp into the protocal buffers timestamp format.
	sts := &tspb.Timestamp{
		Seconds: start,
		Nanos:   0,
	}
	ets := &tspb.Timestamp{
		Seconds: end,
		Nanos:   0,
	}

	// Generate our labels.
	labels := []*pb.Label{}
	labels = append(labels, &pb.Label{Name: "Source", Value: "Command"})

	// Perform our gRPC request.
	res, err := client.SelectBulkTimeSeriesData(ctx, &pb.FilterReq{
		Labels: labels,
		Metric: metric,
		Start: sts,
		End: ets,
	})

	if err != nil {
		log.Fatalf("could not add: %v", err)
	}

	// Print out the gRPC response.
	log.Printf("Server Response:")
	fmt.Println(res)
}

var selectBulkCmd = &cobra.Command{
	Use:   "select",
	Short: "List time-series data",
	Long:  `Connect to the gRPC server and return list of time-series data results based on a selection filter.`,
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
		doSelectRow()
	},
}
