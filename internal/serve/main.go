package serve

import (
	"log/slog"
)

func Run() error {
	slog.Info("Starting dashboard web service")
	slog.Debug("Debugging enabled")

	return nil
}
