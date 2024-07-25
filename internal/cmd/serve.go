package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

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
	// serveConfigFile is a place holder for Cobra to store the alternate
	// configuration file, if set, from the command-line.
	serveConfigFile string

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
	        --endpoint-uri https://development.dashboard.n3t.uk
	  `), "\n"),

		RunE: runServe,
	}
)

// init will initialise the command-line settings for `serveCmd` command,
// including any command-specific flags.
func init() {
	flags := serveCmd.PersistentFlags()
	flags.StringVarP(&serveConfigFile, "config", "c", "", "Path to the configuration file")

	rootCmd.AddCommand(serveCmd)
}

// runServe will run when the serve command is provided to the command-line
// application, initialising the web service, connecting to third-party
// services, and waiting for events to be sent to it for processing. If there
// was an error processing the configuration or initialising the web service or
// its connections, an `error` will be returned.
func runServe(_ *cobra.Command, _ []string) error {
	err := config.Load(serveConfigName, serveConfigFile)
	if err != nil {
		//nolint:revive,stylecheck // new-line is required to break error and usage
		return fmt.Errorf("\n  %w\n", err)
	}

	// As this is a web service, include more information about the release and
	// build environment to make it easier to track and debug changes from logs
	attrs := &map[string]string{
		"name":       Name,
		"version":    Version,
		"commit":     Commit,
		"arch":       Architecture,
		"build-date": BuildDate,
	}

	serve.Prepare()
	logger.Start(attrs)

	err = serve.Run()

	return err
}
