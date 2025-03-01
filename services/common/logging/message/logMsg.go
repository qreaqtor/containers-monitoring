package logmsg

import (
	"log/slog"
)

type LogMsg struct {
	URL    string
	Method string
	Text   string
	Status int
}

// Возвращает структуру, которая пишет логи с помощью logger.
// Остальные поля - информация, которая будет выводиться.
func NewLogMsg(url, method string) *LogMsg {
	return &LogMsg{
		URL:    url,
		Method: method,
	}
}

func (msg *LogMsg) WithText(text string) *LogMsg {
	return &LogMsg{
		Text:   text,
		Status: msg.Status,
		URL:    msg.URL,
		Method: msg.Method,
	}
}

func (msg *LogMsg) WithStatus(status int) *LogMsg {
	return &LogMsg{
		Text:   msg.Text,
		Status: status,
		URL:    msg.URL,
		Method: msg.Method,
	}
}

func (msg *LogMsg) Info() {
	slog.Info(msg.Text, getArgs(msg)...)
}

func (msg *LogMsg) Error() {
	slog.Error(msg.Text, getArgs(msg)...)
}

func getArgs(msg *LogMsg) []any {
	args := []any{
		slog.String("url", msg.URL),
		slog.String("method", msg.Method),
	}

	if msg.Status > 0 {
		args = append(args, slog.Int("status", msg.Status))
	}

	return args
}
