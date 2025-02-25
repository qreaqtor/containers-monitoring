package discard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// always return nil because this handler ignores writting
func (h DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// return h because dont need any attrs to save
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// return h because dont need any groups to save
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// always return false
func (h DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
