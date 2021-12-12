package cmd

import (
	"context"
	"crypto/tls"
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
	autocertDSN  string
	autocertDump string
)

var httpsCmd = &cobra.Command{
	Use:   "https",
	Short: "Start HTTPs server",
	Run:   httpsCommandRun,
}

func init() {
	httpsCmd.PersistentFlags().StringVarP(&caFile, "ca-file", "", "", "file containing the CA")
	httpsCmd.PersistentFlags().StringVarP(&keyFile, "key-file", "", "", "file containing the private key")
	httpsCmd.PersistentFlags().StringVarP(&certFile, "cert-file", "", "", "file containing the public cert(s)")

	httpsCmd.PersistentFlags().BoolVarP(&autoGenerate, "auto-cert", "", false, "generate the certs")
	httpsCmd.PersistentFlags().StringVarP(&autocertDSN, "auto-cert-dsn", "", "localhost", "domain names (separated by ,) for the ssl certs")
	httpsCmd.PersistentFlags().StringVarP(&autocertDump, "auto-cert-rootfile", "", "", "root cert file to dump")

	RootCmd.AddCommand(httpsCmd)
}

func httpsCommandRun(_ *cobra.Command, _ []string) {
	var serverTLS []tls.Certificate
	var rootTLS []byte

	if autoGenerate {
		if server, root, err := certs.GenerateCerts(autocertDSN); err != nil {
			log.Fatal("Cannot generate certs: ", err)
		} else {
			serverTLS = server
			rootTLS = root
		}
	} else {
		if caFile == "" || keyFile == "" || certFile == "" {
			log.Fatal("ca-file, key-file, cert-file are required")
		}

		if server, root, err := certs.LoadCerts(caFile, keyFile, certFile); err != nil {
			log.Fatal("Cannot load certs: ", err)
		} else {
			serverTLS = server
			rootTLS = root
		}
	}

	if autoGenerate && autocertDump != "" {
		f, err := os.OpenFile(autocertDump, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Cannot write root cert file: ", err)
		}
		f.Write(rootTLS)
		f.Close()
	}

	routerConfig := router.GenerateHTTPRoutes(router.RouterConfig{
		Delay:      delay,
		DelayQuery: delayQuery,
		Mirror:     mirrorPayload,
		Cors:       cors,
	})

	srv, err := httpsserver.CreateAndStartHTTPsServer(httpsserver.WebServerConfig{
		Bind:         bindIp,
		Port:         port,
		Certificates: serverTLS,
		RootCA:       rootTLS,
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
