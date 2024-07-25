package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands and
// provides the point from which all other subcommands are attached.
var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [options] command [options]", Name),
	Short: "A server and client for running a web-based dashboard",
	Long: heredoc.Doc(`
	  dashboard is a web-based application which provides a high-level dashboard
	  for receiving and displaying events and related information about any
	  system. For example, when a deployment has been triggered and what is it's
	  final status, or a function has been triggered, or a job has been run.
	`),

	Example: strings.TrimRight(heredoc.Doc(`

	  $ dashboard serve \
	      --endpoint-uri https://development.dashboard.n3t.uk
	  $ dashboard send \
	      --endpoint-uri https://development.dashboard.n3t.uk \
	      --event-id this-is-a-test-message \
	      --status pass \
	      --message 'This is a test message for the dashboard'
	`), "\n"),
}

// init will initialise the command-line settings for `rootCmd` command, as well
// as set up the general application configuration (such as the processing of
// environment variables) and settings.
func init() {
	flags := rootCmd.PersistentFlags()
	flags.StringP("log-level", "l", "info", "Set the logging level (debug, info, warning, error)")

	err := viper.BindPFlag("logging.level", flags.Lookup("log-level"))
	if err != nil {
		panic(err)
	}
}

// Execute executes `rootCmd` and therefore starts the application, with all the
// `init()` functions in each of the files under the `cmd` package helping to
// build up all the configuration options possible to prepare for execution.
// This should only be called by `main.main()`.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
