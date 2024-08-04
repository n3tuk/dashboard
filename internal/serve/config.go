//nolint:mnd // defaults expect hard-coded numberic values in this file
package serve

import (
	"github.com/spf13/viper"
)

const (
	host = "localhost"
	port = 8080
)

// Prepare pre-configured the expected configuration values inside Viper with
// the default values, which can be used through the rest of the code, and can
// be overridden by the YAML files provided to the application.
func Prepare() {
	// Configure the bindings for the web service, setting both the address and
	// port to listen on, and which endpoints should be trusted for processing the
	// X-Forwarded-For headers received from upstream connections
	viper.SetDefault("web.bind.hostname", host)
	viper.SetDefault("web.bind.port", port)
	viper.SetDefault("web.bind.proxies", []string{"127.0.0.1", "::1"})

	// Configure the basic timeouts for the http.Server resource
	viper.SetDefault("web.timeouts.read", 5)
	viper.SetDefault("web.timeouts.write", 10)
	viper.SetDefault("web.timeouts.idle", 30)
	viper.SetDefault("web.timeouts.headers", 2)

	// Provide configuration for the logger, including setting JSON, structured
	// output, and the level of logging output by default
	viper.SetDefault("logging.json", true)
	viper.SetDefault("logging.level", "info")
}
