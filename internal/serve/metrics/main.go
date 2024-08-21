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
	attr   slog.Attr
	router *gin.Engine
	server *http.Server
	health *healthz.Health
}

var ErrServiceNotConfigured = errors.New("service not configured")

func NewService() *Service {
	router := gin.New()

	name := viper.GetString("cluster.name")
	address := viper.GetString("endpoints.bind.address")
	port := viper.GetString("endpoints.bind.port.metrics")

	if viper.GetBool("logging.metrics") {
		router.Use(middleware.Logger())
	}

	router.Use(middleware.Prometheus(name, "metrics"))
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

		health: healthz.NewHealth(),

		attr: slog.Group(
			"cluster",
			slog.String("name", name),
			slog.String("service", "metrics"),
			slog.String("address", address),
			slog.String("port", port),
		),
	}

	Attach(router)
	alive.Attach(router)
	healthz.Attach(router, service.health)

	// Set up the default 404 handler
	router.NoRoute(notFound)

	return service
}

func (s *Service) Start(e chan error) {
	if s.server == nil {
		s.health.Metrics = false
		slog.Error(
			"Failed to start metrics service",
			slog.Group("error", slog.String("message", "service not configured")),
			s.attr,
		)
		e <- ErrServiceNotConfigured
	}

	slog.Info("Starting dashboard metrics service", s.attr)

	s.health.Metrics = true

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.health.Metrics = false
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
	s.health.Terminating = true
}

func (s *Service) SetWebHealth(status bool) {
	s.health.Web = status
}

func (s *Service) SetMetricsHealth(status bool) {
	s.health.Metrics = status
}

func (s *Service) Shutdown(timeout time.Duration) error {
	slog.Info("Shutting down the metrics service", s.attr)

	// Create a context that is used to inform the server it has only a set time
	// to finish the request it is currently handling before being shut down
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.SetMetricsHealth(false)

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
