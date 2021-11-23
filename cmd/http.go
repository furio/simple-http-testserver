package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/furio/simple-http-testserver/lib/http"
	router "github.com/furio/simple-http-testserver/lib/router"

	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start HTTP server",
	Run:   httpCommandRun,
}

func init() {
	RootCmd.AddCommand(httpCmd)
}

func httpCommandRun(_ *cobra.Command, _ []string) {
	routerConfig := router.GenerateHTTPRoutes(router.RouterConfig{
		Delay:      delay,
		DelayQuery: delayQuery,
		Mirror:     mirrorPayload,
		Cors:       cors,
	})

	srv, err := httpserver.CreateAndStartHTTPServer(httpserver.WebServerConfig{
		Bind: bindIp,
		Port: port,
	}, routerConfig)

	if err != nil {
		log.Printf("Cannot start HTTP server due to %s\n", err)
		return
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server (HTTP) forced to shutdown:", err)
	}
}
