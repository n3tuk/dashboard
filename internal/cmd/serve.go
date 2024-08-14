package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/n3tuk/dashboard/internal/config"
	"github.com/n3tuk/dashboard/internal/logger"
	"github.com/n3tuk/dashboard/internal/serve/metrics"
	"github.com/n3tuk/dashboard/internal/serve/web"
)

const (
	// serveConfigName defines the default name of the configuration file for when
	// running the web service for dashboard.
	serveConfigName = "serve.yaml"
)

var (
	// host is the hostname or IPv4/IPv6 address to bind the service to on
	// startup.
	host = "localhost"
	// webPort is the TCP port number to bind the web service to on startup.
	webPort = 8080
	// metricsPort is the TCP port number to bind the metrics service to on startup.
	metricsPort = 8888

	// trustedProxies is a list of IPv4 and/or IPv6 CIDRs which should be trusted
	// for providing the remote Client address.
	trustedProxies = []string{"127.0.0.1", "::1"}

	// headersTimeout is the maximum time to read the full request, after the
	// headers, from the client.
	headersTimeout = 2
	// readTimeout is the maximum time to read the full request, after the
	// headers, from the client.
	readTimeout = 5
	// writeTimeout is the maximum time to write the full response, including
	// the body, to the client.
	writeTimeout = 10
	// idleTimeout is the maximum time to read the headers for the request from
	// the client.
	idleTimeout = 30
	// shutdownTimeout is the maximum time to wait for the web service to finish
	// handling all currently active requests before shutting down.
	shutdownTimeout = 30
	// shutdownMetrics is the maximum time to wait to shut down the metrics
	// service and therefore shut down the application once the web service is
	// closed.
	shutdownMetrics = 5

	// loggerConfig provides the application information which will be used for
	// every log line to help provide context to all logs.
	loggerConfig = &map[string]string{
		"name":       Name,
		"version":    Version,
		"commit":     Commit,
		"arch":       Architecture,
		"build-date": BuildDate,
	}

	// serveCmd represents the serve command for the dashboard application, and will
	// provide the setup and arguments needed for the application to start the web
	// service and start processing events.
	serveCmd = &cobra.Command{
		Use:   "serve [options]",
		Short: "Start the web server to serve dashboard web requests",
		Long: heredoc.Doc(`
		  dashboard serve provides the web service which runs the processing of
		  events submitted to the dashboard, to be saved and pushed out to the
		  clients.
	  `),

		// Add blank line at the top for enforced extra spacing in the output
		Example: strings.TrimRight(heredoc.Doc(`

	    $ dashboard serve --address 0.0.0.0 --web-port 8080 --metrics-port 8081
	  `), "\n"),

		RunE: runServe,
	}
)

// init will initialise the command-line settings for `serveCmd` command,
// including any command-specific flags.
func init() {
	flags := serveCmd.Flags()

	// Flags and default configuration for binding the web service
	viper.SetDefault("endpoints.bind.address", host)
	flags.StringP("address", "a", host, "Address to bind the server to")
	_ = viper.BindPFlag("endpoints.bind.address", flags.Lookup("address"))

	viper.SetDefault("endpoints.bind.port.web", webPort)
	flags.IntP("web-port", "p", webPort, "The port to bind the web service to")
	_ = viper.BindPFlag("endpoints.bind.port.web", flags.Lookup("web-port"))

	viper.SetDefault("endpoints.bind.port.metrics", metricsPort)
	flags.IntP("metrics-port", "m", metricsPort, "The port to bind the metrics service to")
	_ = viper.BindPFlag("endpoints.bind.port.metrics", flags.Lookup("metrics-port"))

	viper.SetDefault("endpoints.proxies", trustedProxies)
	flags.StringSlice("proxies", trustedProxies, "A comma-separated list of CIDRs where trusted proxies are used")
	_ = viper.BindPFlag("endpoints.proxies", flags.Lookup("proxies"))

	// Flags and default configurations for the web service timeouts
	viper.SetDefault("endpoints.timeouts.headers", headersTimeout)
	flags.Int("headers-timeout", headersTimeout, "Timeout (in seconds) to read the headers for the request")
	_ = viper.BindPFlag("endpoints.timeouts.headers", flags.Lookup("headers-timeout"))

	viper.SetDefault("endpoints.timeouts.read", readTimeout)
	flags.Int("read-timeout", readTimeout, "Timeout (in seconds) to read the full request, after the headers")
	_ = viper.BindPFlag("endpoints.timeouts.read", flags.Lookup("read-timeout"))

	viper.SetDefault("endpoints.timeouts.write", writeTimeout)
	flags.Int("write-timeout", writeTimeout, "Timeout (in seconds) to write the full response, including the body")
	_ = viper.BindPFlag("endpoints.timeouts.write", flags.Lookup("write-timeout"))

	viper.SetDefault("endpoints.timeouts.idle", idleTimeout)
	flags.Int("idle-timeout", idleTimeout, "Timeout (in seconds) to keep a connection open between requests")
	_ = viper.BindPFlag("endpoints.timeouts.idle", flags.Lookup("idle-timeout"))

	viper.SetDefault("endpoints.timeouts.shutdown", shutdownTimeout)
	flags.Int("shutdown-timeout", shutdownTimeout, "Timeout (in seconds) to wait for requests to finish")
	_ = viper.BindPFlag("endpoints.timeouts.shutdown", flags.Lookup("shutdown-timeout"))

	rootCmd.AddCommand(serveCmd)
}

// runServe will run when the serve command is provided to the command-line
// application, initialising the web service, connecting to third-party
// services, and waiting for events to be sent to it for processing. If there
// was an error processing the configuration or initialising the web service or
// its connections, an `error` will be returned.
//
//nolint:funlen // ignore
func runServe(_ *cobra.Command, _ []string) error {
	err := config.Load(serveConfigName, configFile)
	if err != nil {
		//nolint:revive,stylecheck // new-line is required to break error and usage
		return fmt.Errorf("\n  %w\n", err)
	}

	gin.SetMode(gin.ReleaseMode)
	logger.Start(loggerConfig)

	// Create a context that listens for the interrupt signal from the Operating
	// System so we can capture it and then trigger a graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	m := metrics.NewService()
	w := web.NewService()

	e := make(chan error)

	// Start the web service first as the metrics service will report the health
	// of the service, so we should be ready to receive requests before the
	// service is reporting as healthy
	go w.Start(e)
	go m.Start(e)

	// Restore default behaviour on the interrupt signal and notify user of shutdown.
	select {
	case <-ctx.Done():
		slog.Info("Shutting down dashboard gracefully")
	case err := <-e:
		slog.Error(
			"Shutting down dashboard due to startup failure",
			slog.Group("error",
				slog.String("message", err.Error()),
			),
		)
	}

	m.PrepareShutdown()

	shutdown := time.Duration(viper.GetInt("endpoints.timeouts.shutdown")) * time.Second
	if err := w.Shutdown(shutdown); err != nil {
		slog.Error(
			"Forced to shut down web service ungracefully",
			slog.Group("error",
				slog.String("message", err.Error()),
			),
		)
	}

	// Only once all the above steps are processed, allow the signals to be
	// processed again, allowing the application to be forcefully terminated, but
	// the client connections have been cleanly closed, so this is acceptable now
	stop()

	shutdown = time.Duration(shutdownMetrics) * time.Second
	if err := m.Shutdown(shutdown); err != nil {
		slog.Error(
			"Forced to shut down metrics service ungracefully",
			slog.Group("error",
				slog.String("message", err.Error()),
			),
		)
	}

	return nil
}
