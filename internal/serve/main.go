package serve

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/n3tuk/dashboard/internal/serve/alive"
	"github.com/n3tuk/dashboard/internal/serve/healthz"

	slogg "github.com/samber/slog-gin"
)

// timeout provides the time allowed to gracefully shut down the service.
const timeout = 30 * time.Second

// Run initiates the setup and startup of the web service, attaching the
// endpoints avoidable in each of the packages for dashboard, as well as
// monitoring for system signals and handling a graceful shutdown of the server
// when called.
func Run() {
	// Create a context that listens for the interrupt signal from the Operating
	// System so we can capture it and then trigger a graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := SetupGin()
	server := SetupServer(router)

	alive.Attach(router)
	healthz.Attach(router)

	// Initialising the server in a goroutine so that it won't block the capture
	// and processing of the system interrupt
	go RunServer(ctx, server)

	// Restore default behaviour on the interrupt signal and notify user of shutdown.
	<-ctx.Done()
	stop()
	slog.InfoContext(ctx,
		"Shutting down dashboard gracefully",
		slog.Group("server",
			slog.String("address", server.Addr),
			slog.String("timeout", timeout.String()),
		),
	)

	// Create a context that is used to inform the server it has only a set time
	// to finish the request it is currently handling before being shut down
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx,
			"Forced to shut down dashboard ungracefully",
			slog.Group("server",
				slog.String("address", server.Addr),
				slog.String("timeout", timeout.String()),
			),
			slog.Group("error",
				slog.String("message", err.Error()),
			),
		)
	}
}

// SetupGin sets up the Gin package, adding in the standard middleware to
// process log entires and recorder from panics and errors which are unhandled
// by individual handlers.
func SetupGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(Logger())
	r.Use(gin.Recovery())

	proxies := viper.GetStringSlice("web.bind.proxies")
	if len(proxies) > 0 {
		err := r.SetTrustedProxies(proxies)
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

	return r
}

// Logger provides a structured logging logger which can be used by Gin using
// the new slog package, allowing for easy processing of log data.
func Logger() gin.HandlerFunc {
	return slogg.NewWithConfig(
		slog.Default().WithGroup("gin"),
		slogg.Config{
			WithRequestID: true,
		},
	)
}

// SetupServer sets up the http package web service, configuring the bindings
// and timeouts for processing requests, as well as attaching the Gin framework.
func SetupServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr: net.JoinHostPort(
			viper.GetString("web.bind.address"),
			viper.GetString("web.bind.port"),
		),

		ReadTimeout:       time.Duration(viper.GetInt("web.timeouts.read")) * time.Second,
		WriteTimeout:      time.Duration(viper.GetInt("web.timeouts.write")) * time.Second,
		IdleTimeout:       time.Duration(viper.GetInt("web.timeouts.idle")) * time.Second,
		ReadHeaderTimeout: time.Duration(viper.GetInt("web.timeouts.header")) * time.Second,

		Handler: router,
	}
}

// RunServer provides a function to be called as a goroutine which starts up and
// runs the web service in a dedicated thread, allowing the main thread to
// handle startup and shutdown independently.
func RunServer(ctx context.Context, server *http.Server) {
	slog.InfoContext(ctx,
		"Starting dashboard",
		slog.Group(
			"server",
			slog.String("address", server.Addr),
		),
	)

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx,
			"Failed to start web service",
			slog.Group("error",
				slog.String("message", err.Error()),
			),
		)
	}
}
