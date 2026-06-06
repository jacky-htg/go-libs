package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type LogLevel string

const (
	LogTypeDebug LogLevel = "debug"
	LogTypeInfo  LogLevel = "info"
	LogTypeWarn  LogLevel = "warn"
	LogTypeError LogLevel = "error"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

type Options struct {
	Level      slog.Level
	AddSource  bool
	FormatJSON bool
	UseTrace   bool
}

type slogLogger struct {
	logger   *slog.Logger
	useTrace bool
}

func InitLogger(opts *Options) Logger {
	if opts == nil {
		opts = &Options{
			Level:      slog.LevelInfo,
			AddSource:  true,
			FormatJSON: true,
			UseTrace:   false,
		}
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: false, // MATIKAN source otomatis, kita handle manual
	}

	var handler slog.Handler
	if opts.FormatJSON {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	return &slogLogger{
		logger:   slog.New(handler),
		useTrace: opts.UseTrace,
	}
}

func SetLogger(log *slog.Logger) Logger {
	return &slogLogger{
		logger:   log,
		useTrace: true,
	}
}

func (l *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, l.enrichArgs(ctx, msg, LogTypeDebug, args)...)
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, l.enrichArgs(ctx, msg, LogTypeInfo, args)...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, l.enrichArgs(ctx, msg, LogTypeWarn, args)...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, l.enrichArgs(ctx, msg, LogTypeError, args)...)
}

// With mengembalikan logger baru dengan args tambahan
func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{
		logger:   l.logger.With(args...),
		useTrace: l.useTrace,
	}
}

// enrichArgs menambahkan trace ID dan source info ke args
func (l *slogLogger) enrichArgs(ctx context.Context, msg string, level LogLevel, args []any) []any {
	result := make([]any, 0, len(args)+4)

	file := "unknown"
	line := 0
	_, file, line, ok := runtime.Caller(2)
	if ok {
		file = shortenFilePath(file)
	}
	result = append(result,
		"source", slog.GroupValue(
			slog.String("file", file),
			slog.Int("line", line),
		),
	)

	if l.useTrace {
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().IsValid() {
			result = append(result,
				"trace_id", span.SpanContext().TraceID().String(),
				"span_id", span.SpanContext().SpanID().String(),
			)
		}
		if level == LogTypeError && span.IsRecording() {
			go additionalErrorSpanAttr(ctx, msg, file, line, args)
		}
	}

	result = append(result, args...)
	return result
}

func additionalErrorSpanAttr(ctx context.Context, msg string, file string, line int, args []any) {
	var detailErrMsg string = msg
	for _, arg := range args {
		if attr, ok := arg.(slog.Attr); ok && attr.Key == "error" {
			detailErrMsg = attr.Value.String()
			break
		}
	}

	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, msg)
	span.SetAttributes(
		attribute.String("error.type", "exception"),
		attribute.String("exception.message", detailErrMsg),
		attribute.String("file", file),
		attribute.Int("line", line),
	)
}

func shortenFilePath(path string) string {
	// Cari pattern yang umum
	patterns := []string{"/pkg/", "/delivery/", "/mapper/", "/repository/", "/service/", "/handler/", "/dto/", "/model/"}
	for _, p := range patterns {
		if idx := strings.Index(path, p); idx != -1 {
			return path[idx+1:]
		}
	}
	// Fallback: ambil 2 komponen terakhir
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return path
}
