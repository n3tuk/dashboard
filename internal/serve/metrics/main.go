package metrics

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/n3tuk/dashboard/internal/serve/metrics/alive"
	"github.com/n3tuk/dashboard/internal/serve/metrics/healthz"
	"github.com/n3tuk/dashboard/internal/serve/middleware"
)

type Service struct {
	attr         slog.Attr
	router       *gin.Engine
	server       *http.Server
	shuttingDown bool
}

func NewService() *Service {
	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Prometheus())
	router.Use(gin.Recovery())

	proxies := viper.GetStringSlice("endpoints.proxies")
	if len(proxies) > 0 {
		err := router.SetTrustedProxies(proxies)
		if err != nil {
			slog.Error(
				"Unable to configure trusted proxies",
				slog.Group(
					"error",
					slog.String("message", err.Error()),
				),
			)
		}
	}

	address := viper.GetString("endpoints.bind.address")
	port := viper.GetString("endpoints.bind.port.metrics")

	service := &Service{
		router: router,
		server: &http.Server{
			ReadTimeout:       time.Duration(viper.GetInt("endpoints.timeouts.read")) * time.Second,
			WriteTimeout:      time.Duration(viper.GetInt("endpoints.timeouts.write")) * time.Second,
			IdleTimeout:       time.Duration(viper.GetInt("endpoints.timeouts.idle")) * time.Second,
			ReadHeaderTimeout: time.Duration(viper.GetInt("endpoints.timeouts.header")) * time.Second,

			Addr:    net.JoinHostPort(address, port),
			Handler: router,
		},

		shuttingDown: false,

		attr: slog.Group(
			"server",
			slog.String("service", "metrics"),
			slog.String("address", address),
			slog.String("port", port),
		),
	}

	Attach(router)
	alive.Attach(router)
	healthz.Attach(router, &service.shuttingDown)

	// Set up the default 404 handler
	router.NoRoute(notFound)

	return service
}

func (s *Service) Start(e chan error) {
	if s.server == nil {
		slog.Error(
			"Failed to start metrics service",
			slog.Group("error", slog.String("message", "service not configured")),
			s.attr,
		)
	}

	slog.Info("Starting dashboard metrics service", s.attr)

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"Failed to start metrics service",
			slog.Group("error", slog.String("message", err.Error())),
			s.attr,
		)
		e <- err
	}
}

func (s *Service) PrepareShutdown() {
	slog.Info("Preparing for metrics service for web service shutdown", s.attr)
	s.shuttingDown = true
}

func (s *Service) Shutdown(timeout time.Duration) error {
	slog.Info("Shutting down the metrics service", s.attr)

	// Create a context that is used to inform the server it has only a set time
	// to finish the request it is currently handling before being shut down
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Attach takes a reference to the Gin engine and attaches all the expected
// endpoints which cam be used by clients through this package.
func Attach(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// notFound provides the default handler for requests to the API which have no
// route, and therefore cannot be processed, necessitating a 404 (Page Not
// Found) response back to the client.
func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "page-not-found",
		"code":    http.StatusNotFound,
		"message": "Page not found",
	})
}
