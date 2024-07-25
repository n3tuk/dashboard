// The `config` package provides the centralised management of the configuration
// for the runtime of the application, by managing the loading of the default
// (or specified) configuration file, and linking in the environment variables.
package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// Paths sets the default paths which will be used to look for the default
// configuration file for each of the application commands.
var Paths = []string{
	"$HOME/.config/dashboard",
	"/etc/dashboard",
}

// Load prepares the environment for processing the configuration file for the
// application command, where `name` is the default filename to be searched for
// in `Paths`, otherwise, if set, `file` will override the search and point, and
// then reads the file into `Viper`.
func Load(name, file string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("dashboard")

	// Use SetEnvKeyReplacer to change any dash or period in the configuration
	// path to an underscore to simplify configuration, and so, for example,
	// `web.log-path` would change from `DASHBOARD_WEB.LOG-PATH` to
	// `DASHBOARD_WEB_LOG_PATH` instead
	viper.SetEnvKeyReplacer(
		strings.NewReplacer(
			// These must be set in pairs for each old/new replacement
			"-", "_",
			".", "_",
		),
	)

	if file != "" {
		viper.SetConfigFile(file)

		if err := read(); err != nil {
			return err
		}

		return nil
	}

	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	for _, path := range Paths {
		viper.AddConfigPath(path)
	}

	err := read()
	if err != nil {
		var expected *NotFoundError

		if errors.As(err, &expected) {
			return nil
		}
	}

	return err
}

// read loads the `file` to configure the dashboard send or serve service.
func read() error {
	err := viper.ReadInConfig()
	if err == nil {
		return nil
	}

	var check viper.ConfigFileNotFoundError
	if errors.As(err, &check) {
		return NewNotFoundError(
			viper.GetViper().ConfigFileUsed(),
			"unable to find the configuration file",
			err,
		)
	}

	return NewLoadError(
		viper.GetViper().ConfigFileUsed(),
		"unable to process configuration file",
		err,
	)
}
