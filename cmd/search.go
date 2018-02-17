// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>

package cmd

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/theshadow/ushadow/audify/service"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search TAGS",
	Short: "Perform a search against the audify.fm API",
	Long: `Makes a request against the audify.fm `,
	Example: `audify "president trump" mars`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing positional argument TAGS")
		}

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
		stream, err := c.Search(
			ctx,
			&pb.SearchRequest{
				Tags: tags,
			},
		)

		if err != nil {
			return fmt.Errorf("unable to make request! %s", err)
		}

		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Printf("%#v\n", in)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
