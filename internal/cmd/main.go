// The `cmd` package provides the command-line interface configuration for the
// dashboard application and configures how it can be called, and as such, what
// it needs to be at runtime.
//
// For example, calling `dashboard send` will focus on sending events to an
// endpoint, while `dashboard serve` will provide a web service endpoint and
// connect to DyamoDB and ActiveMQ to process events received and push them out
// to web clients.
package cmd
