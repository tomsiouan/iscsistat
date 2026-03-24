package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"iscsistat/internal/config"

	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

type Server struct {
	httpServer *http.Server
	cfg        config.Config
	logger     *zap.Logger
}

func loadClientCAs(logger *zap.Logger, caPath string) *x509.CertPool {
	clientCA := x509.NewCertPool()

	caCert, err := os.ReadFile(caPath)
	if err != nil {
		logger.Fatal("Could not load client CA", zap.Error(err), zap.String("path", caPath))
	}

	if ok := clientCA.AppendCertsFromPEM(caCert); !ok {
		logger.Fatal("Failed to append client CA", zap.String("path", caPath))
	}

	return clientCA
}

func NewServer(cfg config.Config, logger *zap.Logger) *Server {
	s := &Server{
		cfg:    cfg,
		logger: logger,
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if cfg.HTTP.TLS.ClientCAPath != "" {
		tlsConfig.ClientCAs = loadClientCAs(logger, s.cfg.HTTP.TLS.ClientCAPath)
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		logger.Info("mTLS security enabled", zap.String("ca", cfg.HTTP.TLS.ClientCAPath))
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      s.routes(),
		TLSConfig:    tlsConfig,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http2.ConfigureServer(srv, &http2.Server{})

	s.httpServer = srv
	return s
}

func (s *Server) Serve() error {
	var err error
	if s.cfg.HTTP.TLS.Enabled {
		s.logger.Info("starting server with TLS", zap.String("addr", s.httpServer.Addr))
		err = s.httpServer.ListenAndServeTLS(s.cfg.HTTP.TLS.CertFile, s.cfg.HTTP.TLS.KeyFile)
	} else {
		s.logger.Info("starting server without TLS", zap.String("addr", s.httpServer.Addr))
		err = s.httpServer.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("server error during ListenAndServe", zap.Error(err))
		return fmt.Errorf("server ListenAndServe failed: %w", err)
	}

	s.logger.Info("server stopped")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
