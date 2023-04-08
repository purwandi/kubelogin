package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/purwandi/kubelogin/etcd"
	"github.com/purwandi/kubelogin/keycloak"
)

type ServerConfig struct {
	Port                string
	CertificateFile     string
	CertiticateKeyFile  string
	KubernetesApiServer string
	Keycloak            keycloak.Config
	Etcd                etcd.EtcdClient
}

type Server struct {
	config ServerConfig
	echo   *echo.Echo
}

func NewServer(cfg ServerConfig) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := Handler{
		k8sHost:  cfg.KubernetesApiServer,
		keycloak: cfg.Keycloak,
		etcd:     cfg.Etcd,
	}

	e.POST("/", handler.Authenticate)
	e.GET("/etcd/metrics", handler.EtcdMetrics)
	e.GET("/ping", handler.Ping)

	return &Server{
		config: cfg,
		echo:   e,
	}
}

func (s *Server) Run() {
	go func() {
		var (
			port string = fmt.Sprintf(":%s", s.config.Port)
			err  error
		)

		if s.config.CertificateFile != "" && s.config.CertiticateKeyFile != "" {
			err = s.echo.StartTLS(port, s.config.CertificateFile, s.config.CertiticateKeyFile)
		} else {
			err = s.echo.Start(port)
		}

		if err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		s.echo.Logger.Fatal(err)
	}
}
