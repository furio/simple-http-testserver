package http

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: Add TLS support
type WebServerConfig struct {
	Bind         string
	Port         int
	Certificates []tls.Certificate
	RootCA       []byte
}

func CreateAndStartHTTPsServer(config WebServerConfig, ginRouter *gin.Engine) (*http.Server, error) {
	rootCAPool := x509.NewCertPool()
	rootCAPool.AppendCertsFromPEM(config.RootCA)

	tlsConfig := tls.Config{
		RootCAs:      rootCAPool,
		Certificates: config.Certificates,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,

			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	}

	srvHTTP := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", config.Bind, config.Port),
		TLSConfig: &tlsConfig,
		Handler:   ginRouter,
	}

	go func() {
		log.Printf("Starting HTTPs server on: %s:%d\n", config.Bind, config.Port)

		if err := srvHTTP.ListenAndServeTLS("", ""); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTPs server error: %s\n", err)
		}
	}()

	return srvHTTP, nil
}
