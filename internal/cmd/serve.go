package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/n3tuk/dashboard/internal/config"
	"github.com/n3tuk/dashboard/internal/logger"
	"github.com/n3tuk/dashboard/internal/serve"
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
	// port is the TCP port number to bind the web service to on startup.
	port = 8080

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

	// serveCmd represents the serve command for the dashboard application, and will
	// provide the setup and arguments needed for the application to start the web
	// service and start processing events.
	serveCmd = &cobra.Command{
		Use:     "serve [options]",
		Aliases: []string{"web"},
		Short:   "Start the web server to serve dashboard web requests",
		Long: heredoc.Doc(`
		  dashboard serve provides the web service which runs the processing of
		  events submitted to the dashboard, to be saved and pushed out to the
		  clients.
	  `),

		// Add blank line at the top for enforced extra spacing in the output
		Example: strings.TrimRight(heredoc.Doc(`

	    $ dashboard serve \
	        --address 0.0.0.0 \
	        --port 8080
	  `), "\n"),

		RunE: runServe,
	}
)

// init will initialise the command-line settings for `serveCmd` command,
// including any command-specific flags.
func init() {
	flags := serveCmd.Flags()

	// Flags and default configuration for binding the web service
	viper.SetDefault("web.bind.address", host)
	flags.StringP("address", "a", host, "Address to bind the server to")
	_ = viper.BindPFlag("web.bind.address", flags.Lookup("address"))

	viper.SetDefault("web.bind.port", port)
	flags.IntP("port", "p", port, "The port to bind the server to")
	_ = viper.BindPFlag("web.bind.port", flags.Lookup("port"))

	viper.SetDefault("web.proxies", trustedProxies)
	flags.StringSlice("proxies", trustedProxies, "A comma-separated list of CIDRs where trusted proxies are used")
	_ = viper.BindPFlag("web.proxies", flags.Lookup("proxies"))

	// Flags and default configurations for the web service timeouts
	viper.SetDefault("web.timeouts.headers", headersTimeout)
	flags.Int("headers-timeout", headersTimeout, "Timeout (in seconds) to read the headers for the request")
	_ = viper.BindPFlag("web.timeouts.headers", flags.Lookup("headers-timeout"))

	viper.SetDefault("web.timeouts.read", readTimeout)
	flags.Int("read-timeout", readTimeout, "Timeout (in seconds) to read the full request, after the headers")
	_ = viper.BindPFlag("web.timeouts.read", flags.Lookup("read-timeout"))

	viper.SetDefault("web.timeouts.write", writeTimeout)
	flags.Int("write-timeout", writeTimeout, "Timeout (in seconds) to write the full response, including the body")
	_ = viper.BindPFlag("web.timeouts.write", flags.Lookup("write-timeout"))

	viper.SetDefault("web.timeouts.idle", idleTimeout)
	flags.Int("idle-timeout", idleTimeout, "Timeout (in seconds) to keep a connection open between requests")
	_ = viper.BindPFlag("web.timeouts.idle", flags.Lookup("idle-timeout"))

	rootCmd.AddCommand(serveCmd)
}

// runServe will run when the serve command is provided to the command-line
// application, initialising the web service, connecting to third-party
// services, and waiting for events to be sent to it for processing. If there
// was an error processing the configuration or initialising the web service or
// its connections, an `error` will be returned.
func runServe(_ *cobra.Command, _ []string) error {
	err := config.Load(serveConfigName, configFile)
	if err != nil {
		//nolint:revive,stylecheck // new-line is required to break error and usage
		return fmt.Errorf("\n  %w\n", err)
	}

	// As this is a web service, include more information about the release and
	// build environment to make it easier to track and debug changes from logs
	logger.Start(&map[string]string{
		"name":       Name,
		"version":    Version,
		"commit":     Commit,
		"arch":       Architecture,
		"build-date": BuildDate,
	})

	serve.Run()

	return nil
}
