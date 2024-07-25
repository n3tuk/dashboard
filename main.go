package main

import (
	"github.com/n3tuk/dashboard/internal/cmd"
)

var (
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
)

func main() {
	cmd.Branch = Branch
	cmd.Commit = Commit
	cmd.Version = Version
	cmd.Architecture = Architecture
	cmd.BuildDate = BuildDate

	cmd.Execute()
}
