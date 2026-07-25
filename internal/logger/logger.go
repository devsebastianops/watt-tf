package logger

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

var logger *slog.Logger

func init() {
	// Erstelle den Charm-Logger und konfiguriere ihn als slog.Logger
	charmLogger := log.New(os.Stdout)
	charmLogger.SetLevel(log.InfoLevel)
	charmLogger.SetReportTimestamp(false) // Lässt das CLI cleaner wirken (optional)

	// Charm bietet direkt einen slog-konpatiblen Handler an!
	logger = slog.New(charmLogger)
}

func SetUp(verbose bool) {
	if verbose {
		SetDebug()
	} else {
		SetLevel(slog.LevelInfo)
	}
}

func Debug(msg string, args ...any) { logger.Debug(msg, args...) }
func Info(msg string, args ...any)  { logger.Info(msg, args...) }
func Warn(msg string, args ...any)  { logger.Warn(msg, args...) }
func Error(msg string, args ...any) { logger.Error(msg, args...) }

func SetLevel(level slog.Level) {
	charmLogger := log.New(os.Stdout)
	charmLogger.SetLevel(slogLevelToCharm(level))
	logger = slog.New(charmLogger)
}

func SetDebug() {
	SetLevel(slog.LevelDebug)
}

func SetJSONHandler() {
	// Fallback oder echter JSON-Modus, falls gewünscht
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger = slog.New(handler)
}

func SetTextHandler() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger = slog.New(handler)
}

// Hilfsfunktion zum Mappen der Log-Levels
func slogLevelToCharm(level slog.Level) log.Level {
	switch {
	case level <= slog.LevelDebug:
		return log.DebugLevel
	case level >= slog.LevelError:
		return log.ErrorLevel
	case level >= slog.LevelWarn:
		return log.WarnLevel
	default:
		return log.InfoLevel
	}
}
