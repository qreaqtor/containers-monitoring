package logmsg

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type LogMsg struct {
	operation uuid.UUID
	URL       string
	Method    string
	Text      string
	Status    int
}

// Возвращает структуру, которая пишет логи с помощью logger.
// Остальные поля - информация, которая будет выводиться.
func NewLogMsg(ctx context.Context, url, method string) *LogMsg {
	return &LogMsg{
		URL:       url,
		Method:    method,
		operation: ExtractOperationID(ctx),
	}
}

func (msg *LogMsg) WithText(text string) *LogMsg {
	return &LogMsg{
		Text:      text,
		Status:    msg.Status,
		operation: msg.operation,
		URL:       msg.URL,
		Method:    msg.Method,
	}
}

func (msg *LogMsg) WithStatus(status int) *LogMsg {
	return &LogMsg{
		Text:      msg.Text,
		Status:    status,
		operation: msg.operation,
		URL:       msg.URL,
		Method:    msg.Method,
	}
}

func (msg *LogMsg) Info() {
	slog.Info(msg.Text, getArgs(msg)...)
}

func (msg *LogMsg) Error() {
	slog.Error(msg.Text, getArgs(msg)...)
}

func getArgs(msg *LogMsg) []any {
	return []any{
		"status", msg.Status,
		"url", msg.URL,
		"method", msg.Method,
		"operation", msg.operation,
	}
}
