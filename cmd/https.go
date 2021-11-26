package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	certs "github.com/furio/simple-http-testserver/lib/certs"
	httpsserver "github.com/furio/simple-http-testserver/lib/https"
	router "github.com/furio/simple-http-testserver/lib/router"

	"github.com/spf13/cobra"
)

var (
	caFile       string
	keyFile      string
	certFile     string
	autoGenerate bool
	certDSN      string
)

var httpsCmd = &cobra.Command{
	Use:   "https",
	Short: "Start HTTPs server",
	Run:   httpCommandRun,
}

func init() {
	httpsCmd.PersistentFlags().StringVarP(&caFile, "ca-file", "", "", "file containing the CA")
	httpsCmd.PersistentFlags().StringVarP(&keyFile, "key-file", "", "", "file containing the private key")
	httpsCmd.PersistentFlags().StringVarP(&certFile, "cert-file", "", "", "file containing the public cert(s)")

	httpsCmd.PersistentFlags().BoolVarP(&delayQuery, "auto-cert", "", false, "generate the certs")
	httpsCmd.PersistentFlags().StringVarP(&certDSN, "auto-cert-dsn", "", "localhost", "domain names (separated by ,) for the ssl certs")

	RootCmd.AddCommand(httpCmd)
}

func httpsCommandRun(_ *cobra.Command, _ []string) {
	if autoGenerate {
		if _, _, err := certs.GenerateCerts(certDSN); err != nil {

		}

	} else {
		if caFile == "" || keyFile == "" || certFile == "" {
			log.Fatal("ca-file, key-file, cert-file are required")
		}

	}

	routerConfig := router.GenerateHTTPRoutes(router.RouterConfig{
		Delay:      delay,
		DelayQuery: delayQuery,
		Mirror:     mirrorPayload,
		Cors:       cors,
	})

	srv, err := httpsserver.CreateAndStartHTTPsServer(httpsserver.WebServerConfig{
		Bind: bindIp,
		Port: port,
	}, routerConfig)

	if err != nil {
		log.Printf("Cannot start HTTPs server due to %s\n", err)
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
