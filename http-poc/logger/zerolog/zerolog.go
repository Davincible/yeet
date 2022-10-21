package zerolog

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"go-micro.dev/v4/logger"
)

type Mode uint8

const (
	Production Mode = iota
	Development
)

type ZeroLogger struct {
	zLog zerolog.Logger
	opts Options
}

func (l *ZeroLogger) Init(opts ...logger.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if hs, ok := l.opts.Context.Value(hooksKey{}).([]zerolog.Hook); ok {
		l.opts.Hooks = hs
	}
	if tf, ok := l.opts.Context.Value(timeFormatKey{}).(string); ok {
		l.opts.TimeFormat = tf
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if useDefault, ok := l.opts.Context.Value(useAsDefaultKey{}).(bool); ok && useDefault {
		l.opts.UseAsDefault = useDefault
	}
	if devMode, ok := l.opts.Context.Value(developmentModeKey{}).(bool); ok && devMode {
		l.opts.Mode = Development
	}
	if prodMode, ok := l.opts.Context.Value(productionModeKey{}).(bool); ok && prodMode {
		l.opts.Mode = Production
	}

	skip := 4
	if l.opts.CallerSkipCount > 0 {
		skip = l.opts.CallerSkipCount
	}
	// RESET
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = nil
	zerolog.CallerSkipFrameCount = skip

	switch l.opts.Mode {
	case Development:
		zerolog.ErrorStackMarshaler = func(err error) interface{} {
			fmt.Println(string(debug.Stack()))
			return nil
		}
		consOut := zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				if len(l.opts.TimeFormat) > 0 {
					w.TimeFormat = l.opts.TimeFormat
				}
				w.Out = l.opts.Out
				w.NoColor = false
			},
		)
		// level = logger.DebugLevel
		l.zLog = zerolog.New(consOut).
			Level(zerolog.DebugLevel).
			With().Timestamp().Stack().Logger()
	default: // Production
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		l.zLog = zerolog.New(l.opts.Out).
			Level(zerolog.InfoLevel).
			With().Timestamp().Stack().Logger()
	}

	// Set log Level if not default
	if l.opts.Level != 100 {
		zerolog.SetGlobalLevel(zerologLevel(l.opts.Level))
		l.zLog = l.zLog.Level(zerologLevel(l.opts.Level))
	}

	// Reporting caller
	if l.opts.ReportCaller {
		l.zLog = l.zLog.With().Caller().Logger()
	}

	// Adding hooks if exist
	for _, hook := range l.opts.Hooks {
		l.zLog = l.zLog.Hook(hook)
	}

	// Setting timeFormat
	if len(l.opts.TimeFormat) > 0 {
		zerolog.TimeFieldFormat = l.opts.TimeFormat
	}

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		l.zLog = l.zLog.With().Fields(l.opts.Fields).Logger()
	}

	// Also set it as zerolog's Default logger
	if l.opts.UseAsDefault {
		zlog.Logger = l.zLog
	}

	return nil
}

func (l *ZeroLogger) Fields(fields map[string]interface{}) logger.Logger {
	l.zLog = l.zLog.With().Fields(fields).Logger()
	return l
}

func (l *ZeroLogger) Log(level logger.Level, args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(zerologLevel(level)).Msg(msg)
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	if level == logger.FatalLevel {
		l.opts.ExitFunc(1)
	}
}

func (l *ZeroLogger) Logf(level logger.Level, format string, args ...interface{}) {
	l.zLog.WithLevel(zerologLevel(level)).Msgf(format, args...)
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	if level == logger.FatalLevel {
		l.opts.ExitFunc(1)
	}
}

func (l *ZeroLogger) String() string {
	return "zerolog"
}

func (l *ZeroLogger) Options() logger.Options {
	// FIXME: How to return full opts?
	return l.opts.Options
}

func ProvideZerologLogger(opts ...logger.Option) *ZeroLogger {
	options := Options{
		Options: logger.Options{
			Level:   100,
			Fields:  make(map[string]interface{}),
			Out:     os.Stderr,
			Context: context.Background(),
		},
		ReportCaller: false,
		UseAsDefault: false,
		Mode:         Production,
		ExitFunc:     os.Exit,
	}

	z := &ZeroLogger{opts: options}
	z.Init(opts...)
	return z
}

// NewLogger builds a new logger based on options.
func NewLogger(opts ...logger.Option) logger.Logger {
	// Default options
	options := Options{
		Options: logger.Options{
			Level:   100,
			Fields:  make(map[string]interface{}),
			Out:     os.Stderr,
			Context: context.Background(),
		},
		ReportCaller: false,
		UseAsDefault: false,
		Mode:         Production,
		ExitFunc:     os.Exit,
	}

	l := &ZeroLogger{opts: options}
	_ = l.Init(opts...)
	return l
}

func zerologLevel(level logger.Level) zerolog.Level {
	switch level {
	case logger.TraceLevel:
		return zerolog.TraceLevel
	case logger.DebugLevel:
		return zerolog.DebugLevel
	case logger.InfoLevel:
		return zerolog.InfoLevel
	case logger.WarnLevel:
		return zerolog.WarnLevel
	case logger.ErrorLevel:
		return zerolog.ErrorLevel
	case logger.FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func ZerologToLoggerLevel(level zerolog.Level) logger.Level {
	switch level {
	case zerolog.TraceLevel:
		return logger.TraceLevel
	case zerolog.DebugLevel:
		return logger.DebugLevel
	case zerolog.InfoLevel:
		return logger.InfoLevel
	case zerolog.WarnLevel:
		return logger.WarnLevel
	case zerolog.ErrorLevel:
		return logger.ErrorLevel
	case zerolog.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}

func (l *ZeroLogger) Error(args ...any) {
	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(zerolog.ErrorLevel).Msg(msg)
}

func (l *ZeroLogger) Errorf(f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	l.zLog.WithLevel(zerolog.ErrorLevel).Msg(msg)
}

func (l *ZeroLogger) Info(args ...any) {
	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *ZeroLogger) Infof(f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	l.zLog.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *ZeroLogger) Warning(args ...any) {
	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(zerolog.WarnLevel).Msg(msg)
}

func (l *ZeroLogger) Warningf(f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	l.zLog.WithLevel(zerolog.WarnLevel).Msg(msg)
}

func (l *ZeroLogger) Fatal(args ...any) {
	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(zerolog.FatalLevel).Msg(msg)
}

func (l *ZeroLogger) Fatalf(f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	l.zLog.WithLevel(zerolog.FatalLevel).Msg(msg)
}

func (l *ZeroLogger) Trace(args ...any) {
	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(zerolog.TraceLevel).Msg(msg)
}

func (l *ZeroLogger) Tracef(f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	l.zLog.WithLevel(zerolog.TraceLevel).Msg(msg)
}

func (l *ZeroLogger) Debug(args ...any) {
	msg := fmt.Sprint(args...)
	l.zLog.WithLevel(zerolog.DebugLevel).Msg(msg)
}

func (l *ZeroLogger) Debugf(f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	l.zLog.WithLevel(zerolog.DebugLevel).Msg(msg)
}
