// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>
//

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/theshadow/audify-rpc/service"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the service",
	Long: `Signals the service to perform a graceful shutdown.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 6000)
		defer cancel()

		conn, err := grpc.DialContext(ctx, rpcHost, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("unable to dial service %s", err)
		}

		c := pb.NewAudifyClient(conn)
		_, err = c.Shutdown(ctx, &pb.ShutdownRequest{})
		if err != nil {
			return fmt.Errorf("unable to execute Shutdown()! %s", err)
		}
		return nil
	},
}

func init() {
	stopCmd.Flags().StringVarP(&rpcHost, "connect", "c", ":50051",
		"host and port to connect to")
	RootCmd.AddCommand(stopCmd)
}
