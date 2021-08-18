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
// go run main.go insert_stream -p=50051 -m="solar_biodigester_temperature_in_degrees" -v=50 -t=1600000000 -a=xxx -b=yyy

func init() {
	// The following are required.
	insertStreamCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to attach to the TSD.")
	insertStreamCmd.MarkFlagRequired("metric")
	insertStreamCmd.Flags().Float64VarP(&value, "value", "v", 0.00, "The value to attach to the TSD.")
	insertStreamCmd.MarkFlagRequired("value")
	insertStreamCmd.Flags().Int64VarP(&tsv, "timestamp", "t", 0, "The timestamp to attach to the TSD.")
	insertStreamCmd.MarkFlagRequired("timestamp")
	insertStreamCmd.Flags().StringVarP(&iAccessToken, "access_token", "a", "", "The JWT access token provided with successful authentication")
	insertStreamCmd.MarkFlagRequired("access_token")
	insertStreamCmd.Flags().StringVarP(&iRefreshToken, "refresh_token", "b", "", "The JWT refresh token provided with successful authentication")
	insertStreamCmd.MarkFlagRequired("refresh_token")

	// The following are optional and will have defaults placed when missing.
	insertStreamCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(insertStreamCmd)
}

func doInsertTimeSeriesData() {
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

	stream, err := client.InsertTimeSeriesData(ctx)
	if err != nil {
		log.Fatalf("%v.InsertTimeSeriesData(_) = _, %v", client, err)
	}

	tsd := &pb.TimeSeriesDatumReq{Labels: labels, Metric: metric, Value: value, Timestamp: ts}

	// DEVELOPERS NOTE:
	// To stream from a client to a server using gRPC, the following documentation
	// will help explain how it works. Please visit it if the code below does
	// not make any sense.
	// https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc-1

	if err := stream.Send(tsd); err != nil {
		log.Fatalf("%v.Send(%v) = %v", stream, tsd, err)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Successfully inserted using streaming")
}

var insertStreamCmd = &cobra.Command{
	Use:   "insert_stream",
	Short: "Insert single datum using streaming",
	Long:  `Connect to the gRPC server and send a time-series datum using the streaming RPC.`,
	Run: func(cmd *cobra.Command, args []string) {
		doInsertTimeSeriesData()
	},
}
