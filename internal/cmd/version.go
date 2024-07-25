package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	// Application is the fixed name of the application and does not change if the
	// name of the executable is changed.
	Application = "dashboard"
	// Name is the current name of the application at runtime, such as if the
	// application is called by symlink or renamed, so can provide improved
	// references in internal documentation on local environments.
	Name = "dashboard"
	// Branch is the name of the branch on which the code existed which built this
	// application to track the testing and development as needed.
	Branch = "main"
	// Commit is the SHA ID of the most current commit of the branch on which the
	// code existed which built this application.
	Commit = "00000000"
	// Version is the most current version tag in the repository, detailing the
	// information about the release of the application, and, if committed or
	// uncommitted code is used during the build since the tag, then the suffix
	// `~dev-{date}-{time}` is added to the `Version`.
	Version = "v0.0.0"
	// Architecture is the name of the architecture this application was built
	// for, and therefore runs on.
	Architecture = "unknown"
	// BuildDate is the date and time this application was built.
	BuildDate = "2024-07-01 00:00:00"

	// versionCmd represents the version command for the application, providing an
	// output of the current version of the release and related information such
	// as branch, commit, and architecture it is built for.
	versionCmd = &cobra.Command{
		Use:     "version [--json]",
		Aliases: []string{"v"},
		Short:   "Show the version and build information for this application",
		Run:     runVersion,
	}
)

// init will initialise the command-line settings for `versionCmd` command,
// including any command-specific flags.
func init() {
	// If possible, get the name of this executable at runtime
	if e, err := os.Executable(); err == nil {
		Name = filepath.Base(e)
	}

	flags := versionCmd.PersistentFlags()

	flags.BoolP("json", "j", false, "Output version in JSON format")

	rootCmd.AddCommand(versionCmd)
}

// runVersion will run when the version command is provided to the command-line
// application, providing the release and build-time information for the
// application.
func runVersion(cmd *cobra.Command, _ []string) {
	var output string

	json, err := cmd.Flags().GetBool("json")
	if err == nil && json {
		// Rather than go down the rabbit whole of marshalling JSON, just create the
		// raw strong which can be parsed further down as an alternate output
		output = heredoc.Doc(`
		  {
			  "application":"%s",
			  "version":"%s",
			  "arch":"%s",
			  "build-date":"%s",
			  "commit":"%s",
			  "branch":"%s"
      }
		`)
	} else {
		output = heredoc.Doc(`
			%s %s
		  built for %s on %s
			commit %s on branch %s
		`)
	}

	//nolint:forbidigo // This is a genuine output to the console
	fmt.Printf(output, Application, Version, Architecture, BuildDate, Commit, Branch)
}
