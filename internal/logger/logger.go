package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

var logger *slog.Logger
var level slog.Level = slog.LevelInfo

const LOG_FORMAT_PRETTY = "pretty"
const LOG_FORMAT_TEXT = "text"
const LOG_FORMAT_JSON = "json"

func Debug(msg string, args ...any) { logger.Debug(msg, args...) }
func Info(msg string, args ...any)  { logger.Info(msg, args...) }
func Warn(msg string, args ...any)  { logger.Warn(msg, args...) }
func Error(msg string, args ...any) { logger.Error(msg, args...) }
func Fatal(msg string, args ...any) { logger.Error(msg, args...); os.Exit(1) }

func Debugf(format string, args ...any) { logger.Debug(fmt.Sprintf(format, args...)) }
func Infof(format string, args ...any)  { logger.Info(fmt.Sprintf(format, args...)) }
func Warnf(format string, args ...any)  { logger.Warn(fmt.Sprintf(format, args...)) }
func Errorf(format string, args ...any) { logger.Error(fmt.Sprintf(format, args...)) }
func Fatalf(format string, args ...any) { logger.Error(fmt.Sprintf(format, args...)); os.Exit(1) }

func init() {
	setPrettyHandler()
}

func SetUp(verbose bool, logFormat string, silent bool) {
	if silent {
		setSilent()
	} else if verbose {
		setDebug()
	}

	switch logFormat {
	case LOG_FORMAT_PRETTY:
		setPrettyHandler()
	case LOG_FORMAT_TEXT:
		setTextHandler()
	case LOG_FORMAT_JSON:
		setJSONHandler()
	default:
		setPrettyHandler()
	}
}

func setDebug() {
	level = slog.LevelDebug
}

func setSilent() {
	level = slog.LevelWarn
}

func setPrettyHandler() {
	charmLogger := log.New(os.Stdout)
	charmLogger.SetLevel(slogLevelToCharm(level))
	charmLogger.SetReportTimestamp(false)
	logger = slog.New(charmLogger)
}

func setJSONHandler() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger = slog.New(handler)
}

func setTextHandler() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger = slog.New(handler)
}

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
