package send

import "log/slog"

func Run() error {
	slog.Info("Sending dashboard event")
	slog.Debug("Debugging enabled")

	return nil
}
