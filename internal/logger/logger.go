package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func init() {
	// Initialize logger with custom emoji handler, INFO level by default (less verbose for CLI)
	handler := NewEmojiHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger = slog.New(handler)
}

func SetUp(verbose bool) {
	if verbose {
		SetDebug()
	} else {
		SetLevel(slog.LevelInfo)
	}
}

// Debug logs a debug-level message with optional attributes
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Info logs an info-level message with optional attributes
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Warn logs a warning-level message with optional attributes
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// Error logs an error-level message with optional attributes
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// SetLevel sets the logger level
func SetLevel(level slog.Level) {
	opts := &slog.HandlerOptions{Level: level}
	handler := NewEmojiHandler(os.Stdout, opts)
	logger = slog.New(handler)
}

// SetJSONHandler switches to JSON-based logging (for piping/automation)
func SetJSONHandler() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger = slog.New(handler)
}

// SetDebug enables debug-level logging (verbose output)
func SetDebug() {
	handler := NewEmojiHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger = slog.New(handler)
}

// SetTextHandler switches to standard text logging
func SetTextHandler() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger = slog.New(handler)
}

// ============================================================================
// Custom Emoji Handler Implementation
// ============================================================================

// EmojiHandler is a custom slog handler that formats logs with emojis
type EmojiHandler struct {
	w     io.Writer
	level slog.Level
}

// NewEmojiHandler creates a new emoji handler
func NewEmojiHandler(w io.Writer, opts *slog.HandlerOptions) *EmojiHandler {
	level := slog.LevelInfo // default
	if opts != nil && opts.Level != nil {
		// Try to get the level value
		if l, ok := opts.Level.(slog.Level); ok {
			level = l
		} else if lvl, ok := opts.Level.(*slog.Level); ok {
			level = *lvl
		}
	}
	return &EmojiHandler{
		w:     w,
		level: level,
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *EmojiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle handles the Record and writes it to the output
func (h *EmojiHandler) Handle(ctx context.Context, record slog.Record) error {
	// Get emoji for level
	emoji := levelToEmoji(record.Level)

	// Build context string from attributes
	var contextParts []string
	record.Attrs(func(attr slog.Attr) bool {
		contextParts = append(contextParts, fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any()))
		return true
	})

	// Format: emoji message (Context: key=value key=value)
	var logLine string
	if len(contextParts) > 0 {
		contextStr := strings.Join(contextParts, " ")
		logLine = fmt.Sprintf("%s\t %s (%s)\n", emoji, record.Message, contextStr)
	} else {
		logLine = fmt.Sprintf("%s\t %s\n", emoji, record.Message)
	}

	_, err := h.w.Write([]byte(logLine))
	return err
}

// WithAttrs returns a new handler with the given attributes
func (h *EmojiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns a new handler with the given group
func (h *EmojiHandler) WithGroup(name string) slog.Handler {
	return h
}

// levelToEmoji converts log level to emoji
func levelToEmoji(level slog.Level) string {
	switch level {
	case slog.LevelError:
		return "❌"
	case slog.LevelWarn:
		return "⚠️"
	case slog.LevelInfo:
		return "ℹ️"
	case slog.LevelDebug:
		return "🐛"
	default:
		return "❓"
	}
}
