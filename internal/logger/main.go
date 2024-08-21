package logger

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// Alongside the standard DEBUG environment variable, this program can be run
// as a GitHub Action to send messages, so also look for the RUNNER_DEBUG
// environment variable if DEBUG is not set, allowing debugging of the
// application during GitHub Workflow runs.
var debugVariables = []string{"DEBUG", "RUNNER_DEBUG"}

// Start sets up the logging for the dashboard application, based on the
// combination of default settings for each command, the provided (or found)
// configuration, and the command-line arguments, all managed through `Cobra` and
// `Viper`. Nothing is returned as this sets the application runtime rather than
// a specific object.
func Start(attrs *map[string]string) {
	var items []any
	var groups []slog.Attr

	if attrs != nil && len(*attrs) > 0 {
		for key, value := range *attrs {
			items = append(items, slog.String(key, value))
		}

		groups = append(groups, slog.Group("application", items...))
	}

	level := getLevel()
	method := getMethod()

	switch method {
	case "text":
		slog.SetDefault(
			slog.New(
				getHandlerText(level, groups),
			),
		)
	case "json":
		slog.SetDefault(
			slog.New(
				getHandlerJSON(level, groups),
			),
		)
	}

	slog.SetLogLoggerLevel(level)
}

// getLevel retrieves the required logging level from either the defaults, the
// configuration file, the environment or the command-line flags (in that
// priority), with special consideration given to `DEBUG` and `RUNNER_DEBUG`
// environment variables.
func getLevel() slog.Level {
	if getFirstSet(debugVariables) {
		return slog.LevelDebug
	}

	switch viper.GetString("logging.level") {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// getMethod retrieves the method to be used to output the logs from either the
// defaults, the configuration file, the environment or the command-line flags
// (in that priority).
func getMethod() string {
	switch viper.GetBool("logging.json") {
	case true:
		return "json"
	default:
		return "text"
	}
}

// getFirstSet iterates over a list of environment variable names and finds the
// first one set, and returns the boolean value for that. It does not go through
// each looking for if any are true, but iterates to see if they are set and
// returning the set and found value.
func getFirstSet(names []string) bool {
	for _, name := range names {
		env, ok := os.LookupEnv(name)
		if ok && env == "true" {
			return true
		} else if ok {
			return false
		}
	}

	return false
}

// getHandlerText will set up the text-based handler used (typically) by the
// send command for outputting logs.
func getHandlerText(level slog.Level, attrs []slog.Attr) slog.Handler {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)

	if len(attrs) > 0 {
		return handler.WithAttrs(attrs)
	}

	return handler
}

// getHandlerJSON will set up the JSON-based structured log handler used
// (typically) by the serve command for outputting logs.
func getHandlerJSON(level slog.Level, attrs []slog.Attr) slog.Handler {
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)

	if len(attrs) > 0 {
		return handler.WithAttrs(attrs)
	}

	return handler
}
