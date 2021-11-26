package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: Add TLS support
type WebServerConfig struct {
	Bind string
	Port int
}

func CreateAndStartHTTPsServer(config WebServerConfig, ginRouter *gin.Engine) (*http.Server, error) {
	// TODO: Add TLS support

	srvHTTP := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Bind, config.Port), // Force only localhost
		Handler: ginRouter,
	}

	go func() {
		log.Printf("Starting HTTP server on: %s:%d\n", config.Bind, config.Port)

		if err := srvHTTP.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %s\n", err)
		}
	}()

	return srvHTTP, nil
}
