package comlog

import (
	"log/slog"
	"os"

	"github.com/qreaqtor/containers-monitoring/common/logging/discard"
	"github.com/qreaqtor/containers-monitoring/common/logging/pretty"
)

// Return default slog.Logger if env has unsupported value.
// Posible env values: local, dev, prod, test.
func SetLogger(env string) {
	var handler slog.Handler

	out := os.Stdout

	switch env {
	case local:
		handler = pretty.NewPrettyHandler(out,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			})
	case dev:
		handler = slog.NewJSONHandler(out,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
	case prod:
		handler = slog.NewJSONHandler(out,
			&slog.HandlerOptions{
				Level: slog.LevelError,
			})
	case test:
		handler = discard.NewDiscardHandler()
	default:
		return
	}

	slog.SetDefault(slog.New(handler))
}
