package web

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/n3tuk/dashboard/internal/serve/middleware"
	"github.com/n3tuk/dashboard/internal/serve/web/ping"
)

type Service struct {
	attr   slog.Attr
	router *gin.Engine
	server *http.Server
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
	port := viper.GetString("endpoints.bind.port.web")

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

		attr: slog.Group(
			"server",
			slog.String("service", "web"),
			slog.String("address", address),
			slog.String("port", port),
		),
	}

	ping.Attach(router)
	router.NoRoute(notFound)

	return service
}

func (s *Service) Start(e chan error) {
	if s.server == nil {
		slog.Error(
			"Failed to start web service",
			slog.Group("error", slog.String("message", "service not configured")),
			s.attr,
		)
	}

	slog.Info("Starting dashboard web service", s.attr)

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(
			"Failed to start web service",
			slog.Group("error", slog.String("message", err.Error())),
			s.attr,
		)
		e <- err
	}
}

func (s *Service) Shutdown(timeout time.Duration) error {
	slog.Info("Shutting down the web service", s.attr)

	// Create a context that is used to inform the server it has only a set time
	// to finish the request it is currently handling before being shut down
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// notFound provides the default handler for requests to the API which have no
// route, and therefore cannot be processed, necessitating a 404 (Page Not
// Found) response back to the client.
func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"status":  "page-not-found",
		"message": "The path requested could not be found",
		"path":    c.Request.URL.Path,
	})
}
