package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/metadata"

	pb "github.com/bartmika/mothership-server/proto"
)

// Ex:
// go run main.go insert_bulk -p=50051 -m="solar_biodigester_temperature_in_degrees" -v=50 -t=1600000000 -a=xxx -b=yyy

func init() {
	// The following are required.
	insertBulkRowsCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to attach to the TSD.")
	insertBulkRowsCmd.MarkFlagRequired("metric")
	insertBulkRowsCmd.Flags().Float64VarP(&value, "value", "v", 0.00, "The value to attach to the TSD.")
	insertBulkRowsCmd.MarkFlagRequired("value")
	insertBulkRowsCmd.Flags().Int64VarP(&tsv, "timestamp", "t", 0, "The timestamp to attach to the TSD.")
	insertBulkRowsCmd.MarkFlagRequired("timestamp")
	insertBulkRowsCmd.Flags().StringVarP(&iAccessToken, "access_token", "a", "", "The JWT access token provided with successful authentication")
	insertBulkRowsCmd.MarkFlagRequired("access_token")
	insertBulkRowsCmd.Flags().StringVarP(&iRefreshToken, "refresh_token", "b", "", "The JWT refresh token provided with successful authentication")
	insertBulkRowsCmd.MarkFlagRequired("refresh_token")

	// The following are optional and will have defaults placed when missing.
	insertBulkRowsCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(insertBulkRowsCmd)
}

func doInsertBulkRow() {
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

    datum := &pb.TimeSeriesDatumReq{Labels: labels, Metric: metric, Value: value, Timestamp: ts}
	arr := []*pb.TimeSeriesDatumReq{}
	arr = append(arr, datum)
	data := &pb.TimeSeriesDataListReq{Data: arr,}

	// Perform our gRPC request.
	_, err = client.InsertTimeSeriesData(ctx, data)
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}

	log.Printf("Successfully inserted")
}

var insertBulkRowsCmd = &cobra.Command{
	Use:   "insert_bulk",
	Short: "Insert a multiple time-series datum",
	Long:  `Connect to the gRPC server and sends multiple time-series datum.`,
	Run: func(cmd *cobra.Command, args []string) {
		doInsertBulkRow()
	},
}
