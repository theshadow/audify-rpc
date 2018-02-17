// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/theshadow/ushadow/audify/service"
)

// versionCmd Will return the build version string of the binary.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the build version",
	Long: `Displays the build version of the binary, will display 'dev-build' when the binary isn't an official release.'`,
	Example: `version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 6)
		defer cancel()

		conn, err := grpc.DialContext(ctx, rpcHost, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("unable to dial service %s", err)
		}

		var tags []*pb.Tag
		for _, a := range args {
			tags = append(tags, &pb.Tag{Tag: a})
		}

		c := pb.NewAudifyClient(conn)
		resp, err := c.Version(
			ctx,
			&pb.VersionRequest{},
		)

		if err != nil {
			return fmt.Errorf("unable to make request! %s", err)
		}

		fmt.Println(resp.Version)
		for _, dep := range resp.Dependencies {
			fmt.Println(dep)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
