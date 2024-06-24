package decorator

import (
	"context"
	"fmt"
	"log/slog" // Import the slog package
	"strings"
)

// ApplyCommandDecorators applies logging and metrics decorators to a CommandHandler
func ApplyCommandDecorators[H any](handler CommandHandler[H], logger *slog.Logger, metricsClient MetricsClient) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: commandMetricsDecorator[H]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}

// CommandHandler is an interface for handling commands
type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

// generateActionName generates a name for the handler action
func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
