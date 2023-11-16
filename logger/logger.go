package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

type L struct {
	name  string
	out   io.Writer
	l     *slog.Logger
	level slog.Level
}

// New initializes a new logger. If w is nil, logs will be sent to stdout.
func New(w io.Writer, name, level string, keyvals ...any) *L {
	if w == nil {
		w = os.Stdout
	}

	noColor := true
	if w, ok := w.(*os.File); ok {
		noColor = !isatty.IsTerminal(w.Fd())
	}

	l := slog.New(tint.NewHandler(w, &tint.Options{
		TimeFormat: time.DateTime,
		Level:      ParseLevel(level),
		NoColor:    noColor,
		AddSource:  true,
	}))
	l = l.With("src", name)

	if len(keyvals) > 0 {
		l = l.With(keyvals...)
	}

	return &L{l: l, out: w, level: ParseLevel(level), name: name}
}

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// New returns a sub-logger with the name appended to the existing logger's source
func (l *L) New(name string) *L {
	var nl L
	nl = *l

	nl = *nl.With("src", l.name+"."+name)

	return &nl
}

// With returns a logger with the keyvals appended to the existing logger
func (l *L) With(keyvals ...any) *L {
	l.l = l.l.With(keyvals...)
	return l
}

// Debug logs a message at the debug level
func (l *L) Debug(msg any, keyvals ...any) {
	if !l.l.Enabled(context.Background(), slog.LevelDebug) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf("%s", msg), pcs[0])
	r.Add(keyvals...)
	l.l.Handler().Handle(context.Background(), r)
}

// Info logs a message at the info level
func (l *L) Info(msg any, keyvals ...any) {
	if !l.l.Enabled(context.Background(), slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf("%s", msg), pcs[0])
	r.Add(keyvals...)
	l.l.Handler().Handle(context.Background(), r)
}

// Warn logs a message at the warning level
func (l *L) Warn(msg any, keyvals ...any) {
	if !l.l.Enabled(context.Background(), slog.LevelWarn) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf("%s", msg), pcs[0])
	r.Add(keyvals...)
	l.l.Handler().Handle(context.Background(), r)
}

// Err logs a message at the error level
func (l *L) Err(msg any, keyvals ...any) {
	if !l.l.Enabled(context.Background(), slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf("%s", msg), pcs[0])
	r.Add(keyvals...)
	l.l.Handler().Handle(context.Background(), r)
}

// Fatal logs a message at the fatal level and also exits the program by calling
// os.Exit
func (l *L) Fatal(msg any, keyvals ...any) {
	if !l.l.Enabled(context.Background(), slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf("%s", msg), pcs[0])
	r.Add(keyvals...)
	l.l.Handler().Handle(context.Background(), r)
	os.Exit(1)
}

// Default returns a default logger implementation
func Default() *L {
	return New(nil, "default", "")
}

func Silence() *L {
	return New(io.Discard, "", "")
}
