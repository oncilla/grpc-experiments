package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	legacy "github.com/oncilla/grpc-experiments/service-alias/proto/hello/v1"
	"github.com/oncilla/grpc-experiments/service-alias/proto/org/project/hello/v1"
)

var out io.Writer

type Pather interface {
	CommandPath() string
}

func main() {
	executable := filepath.Base(os.Args[0])
	cmd := &cobra.Command{
		Use:           executable,
		Short:         "gRPC service alias tester",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
	}

	cmd.AddCommand(
		newServerCmd(),
		newClientCmd(),
	)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func newServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "server",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			srv := grpc.NewServer()

			// Register new service.
			hello.RegisterEchoServiceServer(srv, Server{})
			// Register legacy service.
			legacyDescription := hello.EchoService_serviceDesc()
			legacyDescription.ServiceName = "proto.hello.v1.EchoService"
			srv.RegisterService(&legacyDescription, Server{})

			// Start serving the API.
			listener, err := net.Listen("tcp", "localhost:8082")
			if err != nil {
				return err
			}

			return srv.Serve(listener)
		},
	}
	return cmd
}

func newClientCmd() *cobra.Command {
	var flags struct {
		legacy bool
	}
	cmd := &cobra.Command{
		Use:  "client",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			conn, err := grpc.DialContext(ctx, "localhost:8082", grpc.WithInsecure())
			if err != nil {
				return err
			}

			if flags.legacy {
				fmt.Fprintln(out, "Sending legacy request")
				client := legacy.NewEchoServiceClient(conn)
				rep, err := client.Echo(ctx, &legacy.EchoRequest{Message: "legacy"})
				if err != nil {
					return err
				}
				fmt.Fprintln(out, "Response:", rep.Message)
				return nil
			}

			client := hello.NewEchoServiceClient(conn)
			rep, err := client.Echo(ctx, &hello.EchoRequest{Message: "ping"})
			if err != nil {
				return err
			}
			fmt.Fprintln(out, "Response:", rep.Message)
			return nil
		},
	}
	cmd.Flags().BoolVar(&flags.legacy, "legacy", false, "Use legacy request")
	return cmd
}

type Server struct{}

func (Server) Echo(ctx context.Context, req *hello.EchoRequest) (*hello.EchoResponse, error) {
	fmt.Fprintln(out, "Server:", req.Message)
	return &hello.EchoResponse{
		Message: req.Message,
	}, nil
}
