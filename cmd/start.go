// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>

package cmd

import (
	"os"
	"net"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	log "github.com/Sirupsen/logrus"

	pb "github.com/theshadow/audify-rpc/service"
	api2 "github.com/theshadow/audify-rpc/api"

	"golang.org/x/net/context/ctxhttp"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Audify gRPC service.",
	Long: `Starts the service listening on the specified host and port.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New()
		logger.Level = log.Level(debugLevel)
		logger.Formatter = &log.JSONFormatter{}
		log.SetOutput(os.Stdout)

		lis, err := net.Listen("tcp", hostOn)
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}

		done := make(chan struct{})

		api, err := api2.NewWithDoer(
			apiURL,
			logger,
			api2.Retrying(3, api2.BackingOff(1000, api2.Logging(logger, ctxhttp.Do))))
		if err != nil {
			return err
		}

		ver := pb.Version{Binary:BinaryVersion, Dependencies:strings.Split(BinaryDependencies, ";")}

		srv := grpc.NewServer()
		pb.RegisterAudifyServer(srv, pb.New(ver, srv, api, done))
		reflection.Register(srv)

		go srv.Serve(lis)

		<-done

		srv.GracefulStop()

		return nil
	},
}

func init() {
	startCmd.Flags().IntVarP(&debugLevel, "debug", "d", int(log.WarnLevel), "debug level 0-5")
	startCmd.Flags().StringVarP(&apiURL, "api", "a", defaultAPIUrl, "URL for the Audify.fm API.")
	startCmd.Flags().StringVarP(&hostOn, "listen", "l", ":50051",
		"will start the server listening on this host and port")
	RootCmd.AddCommand(startCmd)
}
