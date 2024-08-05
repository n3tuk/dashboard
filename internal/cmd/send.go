package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/n3tuk/dashboard/internal/config"
	"github.com/n3tuk/dashboard/internal/logger"
	"github.com/n3tuk/dashboard/internal/send"
)

const (
	// sendConfigName defines the default name of the configuration file for when
	// sending events to the dashboard endpoint.
	sendConfigName = "send.yaml"
)

// sendCmd represents the send command for the dashboard application, and will
// provide the setup and arguments needed for the application to build an
// event and send it to the dashboard endpoint for processing.
var sendCmd = &cobra.Command{
	Use:   "send [options]",
	Short: "Send events and updates to the dashboard web service",
	Long: heredoc.Doc(`
		dashboard send provides a mechanism to construct and send an event to the
		web service using either an input file or command-line arguments to build
		and/or override the events.
	`),

	// Add blank line at the top for enforced extra spacing in the output
	Example: strings.TrimRight(heredoc.Doc(`

	  $ dashboard send \
	      --endpoint-uri https://development.dashboard.n3t.uk \
	      --event-id this-is-a-test-message \
	      --status pass \
	      --message 'This is a test message for the dashboard'
	`), "\n"),

	RunE: runSend,
}

// init will initialise the command-line settings for `sendCmd` command,
// including any command-specific flags.
func init() {
	viper.SetDefault("endpoint-uri", "http://localhost:8080")
	rootCmd.AddCommand(sendCmd)
}

// runSend will run when the send command is provided to the command-line
// application, providing the building and sending of an event to the dashboard
// endpoint. If there was an error processing the configuration or the event, an
// `error` will be returned.
func runSend(_ *cobra.Command, _ []string) error {
	err := config.Load(sendConfigName, configFile)
	if err != nil {
		//nolint:revive,stylecheck // new-line is required to break error and usage
		return fmt.Errorf("\n  %w\n", err)
	}

	logger.Start(nil)

	return send.Run()
}
