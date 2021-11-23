package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// verbose       bool
	bindIp        string
	port          int
	delay         int
	delayQuery    bool
	mirrorPayload bool
	cors          bool
)

var RootCmd = &cobra.Command{
	Use:   "simple-http-testserver",
	Short: "a simple http/https server to test response",
}

func init() {
	// RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVarP(&bindIp, "bind-ip", "b", "127.0.0.1", "server ip")
	RootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8000, "server port")
	RootCmd.PersistentFlags().IntVarP(&delay, "delay", "d", 0, "delay (in ms) to answer")
	RootCmd.PersistentFlags().BoolVarP(&delayQuery, "delay-query", "", false, "allow to use delay (in ms) in query string to answer")
	RootCmd.PersistentFlags().BoolVarP(&mirrorPayload, "mirror-response", "", false, "send back as response the request")
	RootCmd.PersistentFlags().BoolVarP(&cors, "cors", "", false, "enable cors header")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
