// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>
//

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// defaultAPIUrl is the default URL for the audify.fm API.
const defaultAPIUrl = "https://api.audify.fm/streams/recent"

// apiURL is the URL to hit for making Audify requests
var apiURL string

// cfgFile Is the configuration file to load
var cfgFile string

// rpcHost gRPC host to connect to
var rpcHost string

// debugLevel defines how verbose the logging messages should be
var debugLevel int

// hostOn defines the IP:Port that the gRPC server will host on
var hostOn string

// Version is the build version of the binary
var BinaryVersion = "dev-build"

// Dependencies is a semi-colon delimited string of the binary's dependencies
// and their version in the format <package> <version> core libs have a
// version of null.
var BinaryDependencies = "dev-build null"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "audify",
	Short: "The worlds best content, to go.",
	Long: `Pick a topic of interest and listen as we serve up short-form (TL;DR) audio to your favorite listening devices.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() { 
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.audify.yaml)")
	RootCmd.PersistentFlags().StringVarP(&rpcHost, "connect", "c", ":50051",
		"gRPC host and port to connect to")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		// Search config in home directory and "." with name ".audify" (without extension).
		viper.SetConfigName(".audify")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
